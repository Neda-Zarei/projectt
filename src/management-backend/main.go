package main

import (
	"flag"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/cmd"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/config"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/pkg/logger"
)

var envfile = flag.String("envfile", "", "path to configuration env file")

func main() {
	flag.Parse()
	if *envfile != "" {
		godotenv.Load(*envfile)
	}
	cfg, err := config.ReadEnv()
	if err != nil {
		panic(err)
	}

	log := logger.NewZapLogger(cfg.DevEnv)
	defer log.Sync()

	if err := cmd.Run(cfg, log); err != nil {
		log.Fatal("application stopped", zap.Error(err))
	}
}
