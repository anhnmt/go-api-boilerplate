package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Params struct {
	fx.In

	Config Config
	DB     *gorm.DB
}

func New(p Params) (*casbin.Enforcer, error) {
	if !p.Config.Migrate {
		gormadapter.TurnOffAutoMigrate(p.DB)
	}

	adapter, err := gormadapter.NewAdapterByDBUseTableName(p.DB, p.Config.Prefix, p.Config.TableName)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(p.Config.Model, adapter)
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}
