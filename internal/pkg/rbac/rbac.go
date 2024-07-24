package rbac

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
	adapter, err := gormadapter.NewAdapterByDBUseTableName(p.DB, "", "")
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("model.conf", adapter)
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}
