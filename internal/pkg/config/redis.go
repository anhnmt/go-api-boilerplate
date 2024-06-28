package config

type Redis struct {
	Address        []string `mapstructure:"address"`
	Password       string   `mapstructure:"password"`
	DB             int      `mapstructure:"db"`
	PoolSize       int      `mapstructure:"pool_size" defaultvalue:"100"`
	MinIdleConns   int      `mapstructure:"min_idle_conns"`
	MaxIdleConns   int      `mapstructure:"max_idle_conns"`
	MaxActiveConns int      `mapstructure:"max_active_conns"`
}
