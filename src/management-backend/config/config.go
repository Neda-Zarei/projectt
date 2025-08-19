package config

type Config struct {
	// DevEnv specifies the environment the application runs in.
	DevEnv          bool                  `json:"devEnv" env:"DEV_ENV,required,notEmpty"`
	DB              DBConfig              `json:"db" envPrefix:"DB_"`
	Server          ServerConfig          `json:"server" envPrefix:"SERVER_"`
	JWT             JWTConfig             `json:"jwt" envPrefix:"JWT_"`
	UserPlanService UserPlanServiceConfig `json:"userPlanService" envPrefix:"USER_PLAN_"`
	Arcaptcha       ArcaptchaConfig       `json:"arcaptcha" envPrefix:"ARCAPTCHA_"`
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

type ServerConfig struct {
	Port         int    `json:"port" env:"PORT,required,notEmpty" envDefault:"8080"`
	Host         string `json:"host" env:"HOST,required,notEmpty" envDefault:"0.0.0.0"`
	ReadTimeout  int    `json:"readTimeout" env:"READ_TIMEOUT" envDefault:"30"`
	WriteTimeout int    `json:"writeTimeout" env:"WRITE_TIMEOUT" envDefault:"30"`
}

type JWTConfig struct {
	Secret     string `json:"secret" env:"SECRET,required,notEmpty"`
	Expiration int    `json:"expiration" env:"EXPIRATION" envDefault:"24"` // in hours
}

type UserPlanServiceConfig struct {
	Host string `json:"host" env:"HOST,required,notEmpty"`
	Port int    `json:"port" env:"PORT,required,notEmpty" envDefault:"50051"`
}

type ArcaptchaConfig struct {
	SiteKey   string `json:"siteKey" env:"SITE_KEY,required,notEmpty"`
	SecretKey string `json:"secretKey" env:"SECRET_KEY,required,notEmpty"`
	VerifyURL string `json:"verifyUrl" env:"VERIFY_URL" envDefault:"https://arcaptcha.ir/verify"`
}
