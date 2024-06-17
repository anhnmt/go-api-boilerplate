package config

type Server struct {
	Pprof Pprof `mapstructure:"pprof"`
	Grpc  Grpc  `mapstructure:"grpc"`
}

type Grpc struct {
	Port        int  `mapstructure:"port" defaultvalue:"5000"`
	LogPayload  bool `mapstructure:"log_payload" defaultvalue:"true"`
	Reflection  bool `mapstructure:"reflection" defaultvalue:"true"`
	HealthCheck bool `mapstructure:"health_check" defaultvalue:"true"`
}

type Pprof struct {
	Enable bool `mapstructure:"enable"`
	Port   int  `mapstructure:"port" defaultvalue:"6060"`
}
