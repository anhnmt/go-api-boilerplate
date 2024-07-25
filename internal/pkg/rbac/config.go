package rbac

type Config struct {
	TableName string `validate:"required" mapstructure:"table_name" defaultvalue:"roles"`
	Prefix    string `mapstructure:"prefix"`
	Model     string `validate:"required" mapstructure:"model" defaultvalue:"rbac_model.conf"`
}
