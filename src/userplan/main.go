package main

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/cmd"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/config"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/pkg/logger"
)

var (
	envfile      = flag.String("envfile", "", "path to configuration env file")
	expirePlans  = flag.Bool("expire-plans", false, "run plan expiration process")
	expiringDays = flag.Int("expiring-days", 7, "days threshold for expiring plans check")
)

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

	// Handle CLI commands
	if *expirePlans {
		if err := cmd.ExpirePlans(cfg, log); err != nil {
			log.Fatal("plan expiration failed", zap.Error(err))
			os.Exit(1)
		}
		return
	}

	// Default: run the server
	if err := cmd.Run(cfg, log); err != nil {
		log.Fatal("application stopped", zap.Error(err))
	}
}
