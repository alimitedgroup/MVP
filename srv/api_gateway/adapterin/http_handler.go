package adapterin

import (
	"context"
	"errors"
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	NumRequests     metric.Int64Counter
	Authentications metric.Int64Counter
)

type HTTPHandler struct {
	Engine             *gin.Engine
	ApiGroup           *gin.RouterGroup
	AuthenticatedGroup *gin.RouterGroup
}

type HttpConfig struct {
	Port uint16
}

func NewListener(lc fx.Lifecycle, addr *net.TCPAddr, logger *zap.Logger) (*net.TCPListener, error) {
	ln, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		logger.Error("Failed to listen for HTTP server", zap.String("address", addr.String()), zap.Error(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ln.Close()
		},
	})

	return ln, nil
}

type HttpParams struct {
	fx.In
	Auth      portin.Auth
	Lifecycle fx.Lifecycle
	Listener  *net.TCPListener
	Logger    *zap.Logger
	Meter     metric.Meter
}

func counter(p HttpParams, name string, options ...metric.Int64CounterOption) metric.Int64Counter {
	ctr, err := p.Meter.Int64Counter(name, options...)
	if err != nil {
		p.Logger.Fatal("Failed to setup OpenTelemetry counter", zap.String("name", name), zap.Error(err))
	}
	return ctr
}

func NewHTTPHandler(p HttpParams) *HTTPHandler {
	logger := p.Logger.Named("gin")
	gin.DebugPrintFunc = func(format string, values ...interface{}) {
		format, isWarning := strings.CutPrefix(format, "[WARNING] ")
		if isWarning {
			logger.Warn(format)
		} else {
			logger.Info(format)
		}
	}
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Info(fmt.Sprintf("%s %s -> %s", httpMethod, absolutePath, handlerName))
	}

	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(func(c *gin.Context) {
		NumRequests.Add(c, 1, metric.WithAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.path", c.Request.URL.Path),
		))
		c.Next()
	})
	api := r.Group("/api/v1")
	authenticated := r.Group("/api/v1")
	authenticated.Use(Authentication(p.Auth, p.Logger))

	srv := &http.Server{Handler: r}
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				err := srv.Serve(p.Listener)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					p.Logger.Error("Failed to start HTTP server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := srv.Shutdown(ctx)
			if err != nil {
				p.Logger.Error("Failed to stop HTTP server", zap.Error(err))
				return err
			}

			return nil
		},
	})

	Authentications = counter(p, "authentications")
	NumRequests = counter(p, "num_requests")

	return &HTTPHandler{r, api, authenticated}
}

func Authentication(b portin.Auth, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verdict := "success"
		defer func() {
			Authentications.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		}()

		auth, found := strings.CutPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		if !found {
			logger.Debug("Authorization header not found, or not containing Bearer token")
			ctx.AbortWithStatusJSON(401, dto.MissingToken())
			verdict = "missing token"
			return
		}
		data, err := b.ValidateToken(auth)
		if err != nil {
			if errors.Is(err, business.ErrorTokenExpired) {
				logger.Debug("User provided an expired token")
				verdict = "expired token"
				ctx.AbortWithStatusJSON(401, dto.ExpiredToken())
			} else {
				logger.Debug("Error validating token", zap.Error(err))
				verdict = "invalid token"
				ctx.AbortWithStatusJSON(401, dto.InvalidToken())
			}
			return
		}

		ctx.Set("user_data", data)
		ctx.Next()
	}
}

func RegisterRoutes(http *HTTPHandler, controllers []Controller) {
	for _, controller := range controllers {
		var group *gin.RouterGroup
		if controller.RequiresAuth() {
			group = http.AuthenticatedGroup
		} else {
			group = http.ApiGroup
		}
		group.Handle(controller.Method(), controller.Pattern(), controller.Handler())
	}
}
