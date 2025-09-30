package config

type CfgPostgre struct {
	Host                string `mapstructure:"POSTGRES_HOST"`
	Port                int    `mapstructure:"POSTGRES_PORT"`
	User                string `mapstructure:"POSTGRES_USER"`
	Password            string `mapstructure:"POSTGRES_PASSWORD"`
	DB                  string `mapstructure:"POSTGRES_DB"`
	SSLMode             string `mapstructure:"POSTGRES_SSLMODE"`
	QueryTimeout        int    `mapstructure:"POSTGRES_QUERY_TIMEOUT"`
	PoolConnectMax      int    `mapstructure:"POSTRGES_POOL_CONNECT_MAX"`
	PoolConnectIdleTime int    `mapstructure:"POSTGRES_POOL_CONNECT_IDLE_TIME"`
	PoolConnectLifeTime int    `mapstructure:"POSTGRES_POOL_CONNECT_LIFE_TIME"`
}
