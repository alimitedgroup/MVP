package adapterin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	NumRequests     metric.Int64Counter
	Authentications metric.Int64Counter
	CounterMap      sync.Map
)

type HttpConfig struct {
	Host string
	Port int
}

func ConfigFromEnv() (*HttpConfig, error) {
	config := &HttpConfig{}

	var (
		ok  bool
		err error
	)

	config.Host, ok = os.LookupEnv("HTTP_HOST")
	if !ok {
		return nil, fmt.Errorf("HTTP_HOST environment variable not set")
	}

	port, ok := os.LookupEnv("HTTP_PORT")
	if !ok {
		return nil, fmt.Errorf("HTTP_PORT environment variable not set")
	}
	config.Port, err = strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("HTTP_PORT environment variable should be a number")
	}

	return config, nil
}

type HTTPHandler struct {
	Engine       *gin.Engine
	Authenticate gin.HandlerFunc
}

func NewListener(lc fx.Lifecycle, cfg *HttpConfig, logger *zap.Logger) (*net.TCPListener, error) {
	addrStr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		logger.Fatal(
			"Failed to bind to TCP address",
			zap.Error(err),
			zap.String("addr", addrStr),
		)
	}

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

	observability.CounterSetup(&p.Meter, p.Logger, &Authentications, &CounterMap, "num_api_gateway_authentications")
	observability.CounterSetup(&p.Meter, p.Logger, &NumRequests, &CounterMap, "num_api_gateway_total_requests")

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

	authenticated := Authentication(p.Auth, p.Logger)

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

	return &HTTPHandler{r, authenticated}
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

func CheckRole(roles []types.UserRole, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_data, exist := ctx.Get("user_data")
		if !exist {
			logger.Debug("User data not found in context")
			ctx.AbortWithStatusJSON(401, dto.MissingToken())
			return
		}

		userData, ok := user_data.(types.UserData)
		if !ok {
			logger.Debug("User data in context is not of type UserData")
			ctx.AbortWithStatusJSON(401, dto.MissingToken())
			return
		}

		for _, role := range roles {
			// check if the user has the required role or if the role is RoleNone
			// RoleNone is used to allow access to all users, regardless of their role
			if role == userData.Role || role == types.RoleNone {
				ctx.Next()
				return
			}
		}
		logger.Debug("User does not have the required role", zap.String("user_role", userData.Role.String()))
		ctx.AbortWithStatusJSON(403, dto.AuthFailed())
	}
}

func RegisterRoutes(http *HTTPHandler, logger *zap.Logger, controllers []Controller) {
	for _, controller := range controllers {
		var group *gin.RouterGroup
		if controller.RequiresAuth() {
			group = http.Engine.Group("/api/v1")
			group.Use(http.Authenticate)
			group.Use(CheckRole(controller.AllowedRoles(), logger))
		} else {
			group = http.Engine.Group("/api/v1")
		}
		group.Handle(controller.Method(), controller.Pattern(), controller.Handler())
	}
}
