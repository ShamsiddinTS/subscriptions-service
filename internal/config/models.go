package config

type Settings struct {
	Endpoint Endpoint `mapstructure:"endpoint"`
	Db       Db       `mapstructure:"db"`
}

type Endpoint struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	WriteTimeout int    `mapstructure:"writeTimeout"`
	ReadTimeout  int    `mapstructure:"readTimeout"`
}

type Db struct {
	Host          string
	Port          string
	User          string
	Password      string
	DbName        string `mapstructure:"name"`
	SslMode       string `mapstructure:"sslMode"`
	MigrationsDir string `mapstructure:"migrationsDir"`
}
