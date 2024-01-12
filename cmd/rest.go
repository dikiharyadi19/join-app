package cmd

import (
	"github/yogabagas/join-app/config"
	"github/yogabagas/join-app/transport/rest"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use: "api-serve",
	PreRun: func(cmd *cobra.Command, args []string) {

		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("can't load env", err)
			os.Exit(1)
		}

		config.LoadConfig(configURL)

		sqlDB, _ = InitSQLModule()
		redisClient, _ = InitCache()
	},
	Run: func(cmd *cobra.Command, args []string) {

		rest := rest.NewRest(
			&rest.Option{
				Port:         config.GlobalCfg.App.Port,
				ReadTimeout:  time.Duration(config.GlobalCfg.App.ReadTimeout * int(time.Second)),
				WriteTimeout: time.Duration(config.GlobalCfg.App.WriteTimeout * int(time.Second)),
				Sql:          sqlDB.MySQL,
				Redis:        redisClient.Client,
			},
		)
		go rest.Serve()
		rest.SignalCheck()
	},
}
