package otel

type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	Type     string `validate:"required_with=Endpoint,oneof=grpc http stdout" mapstructure:"format" defaultvalue:"grpc"`
}
