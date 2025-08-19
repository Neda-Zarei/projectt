package config

type Config struct {
	// DevEnv specifies the environment the application runs in.
	DevEnv bool       `json:"devEnv" env:"DEV_ENV,required,notEmpty"`
	DB     DBConfig   `json:"db" envPrefix:"DB_"`
	GRPC   GRPCConfig `json:"grpc" envPrefix:"GRPC_"`
}

type DBConfig struct {
	Host     string `json:"host" env:"HOST,required,notEmpty"`
	Port     uint   `json:"port" env:"PORT,required,notEmpty"`
	DBName   string `json:"dbName" env:"NAME,required,notEmpty"`
	Schema   string `json:"schema" env:"SCHEMA,required,notEmpty"`
	User     string `json:"user" env:"USER,required,notEmpty"`
	Password string `json:"password" env:"PASSWORD,required,notEmpty"`
	AppName  string `json:"appName" env:"APP_NAME,required,notEmpty"`
}

type GRPCConfig struct {
	Port     uint   `json:"port" env:"PORT,required,notEmpty"`
	TLS      bool   `json:"tls" env:"TLS,required,notEmpty"`
	CertFile string `json:"certFile" env:"CERT_FILE,required"`
	KeyFile  string `json:"keyFile" env:"KEY_FILE,required"`
}
