package config

type Config struct {
	DB *DBConfig
	AssetDir string
}

type DBConfig struct {
	Dialect  string
	Username string
	Password string
	Name     string
	Charset  string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "mysql",
			Username: "store",
			Password: "sipuserix",
			Name:     "warung",
			Charset:  "utf8",
		},
		AssetDir: "images",
	}
}
