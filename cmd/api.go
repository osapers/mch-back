package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/osapers/mch-back/internal/provider/kwe"
	"github.com/osapers/mch-back/internal/service/tag"

	"github.com/joho/godotenv"
	"github.com/osapers/mch-back/internal/controller/httpapi"
	"github.com/osapers/mch-back/internal/provider/postgres"
	"github.com/osapers/mch-back/internal/service/event"
	"github.com/osapers/mch-back/internal/service/user"
)

func Launch() error {
	parseEnvFileIfExists(".env")

	pgConn, err := postgres.New()
	if err != nil {
		return fmt.Errorf("cannot connect to db - %s", err)
	}

	defer func() {
		_ = pgConn.Destroy()
	}()

	kweClient, err := kwe.NewClient()
	if err != nil {
		return fmt.Errorf("cannot connect to kwe - %s", err)
	}

	tagService := tag.NewService(pgConn)
	eventService := event.NewService(pgConn)
	userService := user.NewService(pgConn, kweClient)

	srv := httpapi.NewServer(eventService, userService, tagService)
	srv.Start()

	return srv.Shutdown()
}

func parseEnvFileIfExists(filename string) {
	if _, err := os.Stat(filename); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	}
}
