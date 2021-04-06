package casbin

import (
	"github.com/casbin/casbin/v2"
	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"{{ .Extra.pkgpath }}/pkg/auth"
)

// Options is  configuration of database
type Options struct {
	Policy string `yaml:"policy"`
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var err error
	o := new(Options)
	if err = v.UnmarshalKey("casbin", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal db option error")
	}

	logger.Info("load casbin options success", zap.String("policy", o.Policy))

	return o, err
}

// Init 初始化数据库
func New(db *gorm.DB, o *Options) (auth.Enforcer, error) {
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	// enforcer := casbin.NewEnforcer(path, adapter)
	enforcer, err := casbin.NewSyncedEnforcer(o.Policy, adapter)
	if err != nil {
		return nil, err
	}
	enforcer.EnableAutoSave(true)
	// enforcer.StartAutoLoadPolicy(1 * time.Minute)
	return enforcer, nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
