package adapterin

import (
	"context"
	"errors"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

type HTTPHandler struct {
	Engine             *gin.Engine
	ApiGroup           *gin.RouterGroup
	AuthenticatedGroup *gin.RouterGroup
}

type HttpConfig struct {
	Port uint16
}

func NewListener(lc fx.Lifecycle, addr *net.TCPAddr) (*net.TCPListener, error) {
	ln, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		slog.Error("Failed to listen for HTTP server", "address", addr, "error", err)
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ln.Close()
		},
	})

	return ln, nil
}

func NewHTTPHandler(b portin.Auth, lc fx.Lifecycle, ln *net.TCPListener) *HTTPHandler {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(gin.Recovery())
	api := r.Group("/api/v1")
	authenticated := r.Group("/api/v1")
	authenticated.Use(Authentication(b))

	srv := &http.Server{Handler: r}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				err := srv.Serve(ln)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					slog.Error("Failed to start HTTP server", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := srv.Shutdown(ctx)
			if err != nil {
				slog.Error("Failed to stop HTTP server", "error", err)
				return err
			}

			return nil
		},
	})

	return &HTTPHandler{r, api, authenticated}
}

func Authentication(b portin.Auth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth, found := strings.CutPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		if !found {
			ctx.AbortWithStatusJSON(401, dto.MissingToken())
			return
		}
		data, err := b.ValidateToken(auth)
		if err != nil {
			if errors.Is(err, business.ErrorTokenExpired) {
				ctx.AbortWithStatusJSON(401, dto.ExpiredToken())
			} else {
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
