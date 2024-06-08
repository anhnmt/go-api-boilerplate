package gen

import (
	"context"

	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/anhnmt/go-api-boilerplate/db/postgres/query"
	"github.com/anhnmt/go-api-boilerplate/internal/domain/entities"
)

func New(ctx context.Context, db *gorm.DB) {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/infrastructure/persistence/postgresql",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(db.WithContext(ctx)) // reuse your gorm db

	// Generate basic type-safe DAO API
	g.ApplyBasic(
		entities.User{},
	)

	// Generate Type Safe API with Dynamic SQL defined on Query interface
	g.ApplyInterface(func(query.UserQuery) {}, entities.User{})

	// Generate the code
	g.Execute()
}
