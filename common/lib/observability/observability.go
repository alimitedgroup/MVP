package observability

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	prettyconsole "github.com/thessem/zap-prettyconsole"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getBuildInfo() (string, string) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("ReadBuildInfo failed")
		return "unknown", "unknown"
	}

	name := bi.Path

	version := "unknown"

	gitHashRegex := regexp.MustCompile(`^[a-f0-9]{40}$`)

	dirty := false
	for _, kv := range bi.Settings {
		switch kv.Key {
		case "vcs.revision":
			version = kv.Value
			if gitHashRegex.MatchString(version) {
				version = version[:7]
			}
		case "vcs.modified":
			dirty = kv.Value == "true"
		}
	}
	if dirty {
		version = version + "-dirty"
	}

	return name, version
}

// getLogLevel ritorna il livello massimo di log che verranno stampati.
// Se la variabile d'ambiente `LOG_LEVEL` è presente, e ha
// valori "error", "warn", "info" e "debug", allora questa
// funzione ritorna il livello corrispondente.
// Se la variabile d'ambiente non è definita,
// allora viene ritornato il livello di default "debug".
func getLogLevel() zapcore.Level {
	env, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		return zap.DebugLevel
	}

	switch strings.ToLower(env) {
	case "error":
		return zap.ErrorLevel
	case "warn":
		return zap.WarnLevel
	case "info":
		return zap.InfoLevel
	case "debug":
		return zap.DebugLevel
	default:
		panic(fmt.Sprintf("invalid LOG_LEVEL: %s", env))
	}
}

func setupZap(level zapcore.Level) *zap.Logger {
	otelemetry := otelzap.NewCore("core_name")

	stderr := zapcore.NewCore(
		prettyconsole.NewEncoder(prettyconsole.NewEncoderConfig()),
		zapcore.Lock(os.Stderr),
		level,
	)

	return zap.New(zapcore.NewTee(stderr, otelemetry))
}

func setupOtel(otlpUrl string, tempLogger *zap.Logger) func(context.Context) error {
	ctx := context.Background()

	// Setup otel semantic conventions
	name, version := getBuildInfo()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(name),
			semconv.ServiceVersion(version),
			semconv.ServiceInstanceID(uuid.New().String()),
		),
	)
	if err != nil {
		tempLogger.Fatal("Failed to create OpenTelemetry resource", zap.Error(err))
	}

	// gRPC connection to the collector
	conn, err := grpc.NewClient(
		otlpUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		tempLogger.Fatal(
			"Failed to connect to OTLP endpoint",
			zap.String("endpoint", otlpUrl),
			zap.Error(err),
		)
	}

	// Logs
	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		tempLogger.Fatal("Failed to create log exporter", zap.Error(err))
	}
	logProcessor := log.NewBatchProcessor(logExporter)
	logProvider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(logProcessor),
	)
	global.SetLoggerProvider(logProvider)

	// Metrics
	metricsExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		tempLogger.Fatal("Failed to create metric exporter", zap.Error(err))
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricsExporter, sdkmetric.WithInterval(15*time.Second))),
	)
	otel.SetMeterProvider(meterProvider)

	return func(ctx context.Context) (err error) {
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			err = logProvider.Shutdown(ctx)
			wg.Done()
		}()
		go func() {
			err = meterProvider.Shutdown(ctx)
			wg.Done()
		}()
		return
	}
}

// WrapLogger decora un logger con il nome specificato
// Esempio d'uso:
//
//	fx.Decorate(WrapLogger("fancymodule"))
func WrapLogger(name string) func(*zap.Logger) *zap.Logger {
	return func(logger *zap.Logger) *zap.Logger {
		return logger.Named(name)
	}
}

func New(lc fx.Lifecycle) (*zap.Logger, metric.Meter) {
	level := getLogLevel()
	tempLogger := prettyconsole.NewLogger(level)

	otlpUrl, ok := os.LookupEnv("OTLP_URL")
	if ok {
		otelcancel := setupOtel(otlpUrl, tempLogger)
		lc.Append(fx.Hook{OnStop: otelcancel})
	} else {
		tempLogger.Warn("OTLP_URL not set, OpenTelemetry will be disabled")
	}

	logger := setupZap(level)
	lc.Append(fx.Hook{OnStop: func(ctx context.Context) error { return logger.Sync() }})

	name, _ := getBuildInfo()
	return logger, otel.Meter(name)
}

func CounterSetup(meter *metric.Meter, logger *zap.Logger, counter *metric.Int64Counter, counterMap *sync.Map, name string, options ...metric.Int64CounterOption) {
	ctr, err := (*meter).Int64Counter(name, options...)
	if err != nil {
		logger.Fatal("Failed to setup OpenTelemetry counter", zap.String("name", name), zap.Error(err))
	}
	_, loaded := counterMap.LoadOrStore(name, ctr)
	if !loaded {
		*counter = ctr
	}
}

var Module = fx.Options(
	fx.NopLogger,
	fx.Provide(New),
)

var ModuleTest = fx.Options(fx.Provide(TestLogger, TestMeter))

func TestLogger(t *testing.T) *zap.Logger {
	return zaptest.NewLogger(t)
}

func TestMeter() metric.Meter {
	return otel.Meter("test")
}
