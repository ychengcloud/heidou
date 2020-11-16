package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"{{ . }}/pkg/transports/http/middlewares/ginprom"
	netutil "{{ . }}/pkg/utils"
)

type Options struct {
	Host         string
	Port         int
	Mode         string
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

type Server struct {
	o          *Options
	app        string
	host       string
	port       int
	logger     *zap.Logger
	router     *gin.Engine
	httpServer http.Server
}

func NewOptions(v *viper.Viper) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)

	if err = v.UnmarshalKey("http", o); err != nil {
		return nil, err
	}

	return o, err
}

type BaseInitControllers func(r *gin.Engine)
type InitControllers func(r *gin.Engine)

// func NewRouter(o *Options, logger *zap.Logger, init InitControllers, tracer opentracing.Tracer) *gin.Engine {
func NewRouter(o *Options, logger *zap.Logger, baseInit BaseInitControllers, init InitControllers) *gin.Engine {

	// 配置gin
	gin.SetMode(o.Mode)
	r := gin.New()

	r.Use(gin.Recovery()) // panic之后自动恢复
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(cors.New(cors.Config{
		AllowOrigins: o.AllowOrigins,
		AllowMethods: o.AllowMethods,
		AllowHeaders: o.AllowHeaders,
		// 	AllowCredentials: true,
		// 	MaxAge:           12 * time.Hour,
	}))
	r.Use(ginprom.New(r).Middleware()) // 添加prometheus 监控
	// r.Use(ginhttp.Middleware(tracer))

	r.Static("doc", "./doc")
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	pprof.Register(r)

	baseInit(r)
	init(r)

	return r
}

func New(o *Options, logger *zap.Logger, router *gin.Engine) (*Server, error) {
	var s = &Server{
		logger: logger.With(zap.String("type", "http.Server")),
		router: router,
		o:      o,
	}

	return s, nil
}

func (s *Server) Application(name string) {
	s.app = name
}

func (s *Server) Start() error {
	s.host = s.o.Host
	s.port = s.o.Port
	if s.port == 0 {
		s.port = netutil.GetAvailablePort()
	}

	if s.host == "" {
		s.host = netutil.GetLocalIP4()
		if s.host == "" {
			return errors.New("get local ipv4 error")
		}

	}

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	s.httpServer = http.Server{Addr: addr, Handler: s.router}

	s.logger.Info("http server starting ...", zap.String("addr", addr))
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("start http server err", zap.Error(err))
			return
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("stopping http server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // 平滑关闭,等待5秒钟处理
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
	}

	return nil
}
