package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"gorm.io/plugin/opentelemetry/tracing"
)

type Params struct {
	fx.In

	Ctx    context.Context
	Config Config
}

type Postgres struct {
	*gorm.DB

	cfg Config
}

func New(lc fx.Lifecycle, p Params) (*Postgres, error) {
	db, err := parseDBWriter(p.Config.Writer)
	if err != nil {
		return nil, err
	}

	if p.Config.Debug {
		db = db.Debug()
	}

	pg := &Postgres{
		DB:  db,
		cfg: p.Config,
	}

	sqlDB, err := pg.SqlDB()
	if err != nil {
		return nil, err
	}

	newCtx, cancel := context.WithTimeout(p.Ctx, 10*time.Second)
	defer cancel()

	if err = sqlDB.PingContext(newCtx); err != nil {
		return nil, err
	}

	// handle db reader
	if p.Config.Reader.Enable {
		if err = parseDBReader(p.Config.Reader, db); err != nil {
			return nil, err
		}
	}

	if p.Config.Otel {
		if err = db.Use(tracing.NewPlugin()); err != nil {
			return nil, err
		}
	}

	lc.Append(fx.StopHook(pg.Close))

	return pg, nil
}

func (p *Postgres) SqlDB() (*sql.DB, error) {
	return p.DB.DB()
}

func (p *Postgres) Close() error {
	sqlDB, err := p.SqlDB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (p *Postgres) AutoMigrate(dst ...any) error {
	if p.cfg.Migrate {
		return p.DB.AutoMigrate(dst...)
	}

	return nil
}

func parseDBWriter(cfg Base) (*gorm.DB, error) {
	dsn := cfg.ParseDSN()

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}),
		// TranslateError:                           true,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
		CreateBatchSize:                          100,
		NowFunc:                                  time.Now().UTC,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// Config maxConnIdleTime
	maxConnIdleTime, err := time.ParseDuration(cfg.MaxConnIdleTime)
	if err != nil {
		return nil, fmt.Errorf("maxConnIdleTime writer: %w", err)
	}

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxIdleTime(maxConnIdleTime)

	// Config maxConnLifetime
	maxConnLifetime, err := time.ParseDuration(cfg.MaxConnLifetime)
	if err != nil {
		return nil, fmt.Errorf("maxConnLifetime writer: %w", err)
	}

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(maxConnLifetime)

	return db, nil
}

func parseDBReader(cfg Base, db *gorm.DB) error {
	dsn := cfg.ParseDSN()

	// Config maxConnIdleTime
	maxConnIdleTime, err := time.ParseDuration(cfg.MaxConnIdleTime)
	if err != nil {
		return fmt.Errorf("maxConnIdleTime reader: %w", err)
	}

	// Config maxConnLifetime
	maxConnLifetime, err := time.ParseDuration(cfg.MaxConnLifetime)
	if err != nil {
		return fmt.Errorf("maxConnLifetime reader: %w", err)
	}

	return db.Use(
		dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{postgres.Open(dsn.String())},
			// sources/replicas load balancing policy
			Policy: dbresolver.RandomPolicy{},
			// print sources/replicas mode in logger
			TraceResolverMode: true,
		}).
			SetConnMaxIdleTime(maxConnIdleTime).
			SetConnMaxLifetime(maxConnLifetime).
			SetMaxIdleConns(cfg.MaxIdleConns).
			SetMaxOpenConns(cfg.MaxOpenConns),
	)
}
