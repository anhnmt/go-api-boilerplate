package gormgen

import (
	"context"

	"go.uber.org/fx"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/anhnmt/go-api-boilerplate/cmd/gorm-gen/generator"
	"github.com/anhnmt/go-api-boilerplate/internal/model"
)

type Params struct {
	fx.In

	Ctx context.Context
	DB  *gorm.DB
}

func New(shutdowner fx.Shutdowner, p Params) error {
	// Generate code
	g := gen.NewGenerator(gen.Config{
		OutPath: "./gen/gormgen",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(p.DB.WithContext(p.Ctx)) // reuse your gorm db

	// Generate basic type-safe DAO API
	g.ApplyBasic(
		model.User{},
		model.Session{},
	)

	// Generate Type Safe API with Dynamic SQL defined on Query interface
	g.ApplyInterface(func(generator.User) {}, model.User{})
	g.ApplyInterface(func(generator.Session) {}, model.Session{})

	g.Execute()
	return shutdowner.Shutdown()
}
