module {{ . }}

go 1.15

require (
	github.com/HdrHistogram/hdrhistogram-go v1.0.0 // indirect
	github.com/casbin/casbin/v2 v2.17.0
	github.com/casbin/gorm-adapter/v3 v3.0.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-contrib/zap v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1
	github.com/google/wire v0.4.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/spf13/viper v1.7.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.5
)
