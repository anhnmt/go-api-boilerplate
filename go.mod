module github.com/anhnmt/go-api-boilerplate

go 1.22

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240717164558-a6c49f84cc0f.2
	connectrpc.com/vanguard v0.2.0
	github.com/anhnmt/go-fingerprint v1.0.2
	github.com/bufbuild/protovalidate-go v0.6.3
	github.com/bytedance/sonic v1.11.9
	github.com/casbin/casbin/v2 v2.98.0
	github.com/casbin/gorm-adapter/v3 v3.26.0
	github.com/go-playground/validator/v10 v10.22.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0
	github.com/jackc/pgx/v5 v5.6.0
	github.com/redis/go-redis/extra/redisotel/v9 v9.5.3
	github.com/redis/go-redis/v9 v9.6.1
	github.com/rs/cors v1.11.0
	github.com/rs/zerolog v1.33.0
	github.com/samber/lo v1.46.0
	github.com/spf13/cast v1.6.0
	github.com/spf13/viper v1.19.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.53.0
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.28.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/trace v1.28.0
	go.uber.org/automaxprocs v1.5.3
	go.uber.org/fx v1.22.1
	golang.org/x/crypto v0.25.0
	golang.org/x/net v0.27.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gorm.io/driver/postgres v1.5.9
	gorm.io/gen v0.3.26
	gorm.io/gorm v1.25.11
	gorm.io/plugin/dbresolver v1.5.2
	gorm.io/plugin/opentelemetry v0.1.4
)

require (
	connectrpc.com/connect v1.16.2 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/casbin/govaluate v1.2.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/glebarez/go-sqlite v1.22.0 // indirect
	github.com/glebarez/sqlite v1.11.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/google/cel-go v0.20.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.21.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/microsoft/go-mssqldb v1.7.2 // indirect
	github.com/mileusna/useragent v1.3.4 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/redis/go-redis/extra/rediscmd/v9 v9.5.3 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/yihleego/murmurhash3 v0.0.0-20220914065222-8cd2aa986a9d // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240725223205-93522f1f2a9f // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240725223205-93522f1f2a9f // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/datatypes v1.2.1 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/driver/sqlserver v1.5.3 // indirect
	gorm.io/hints v1.1.2 // indirect
	modernc.org/libc v1.55.5 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
	modernc.org/sqlite v1.31.1 // indirect
)
