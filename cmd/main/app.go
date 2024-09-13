package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-rest-api/internal/config"
	"go-rest-api/internal/user"
	"go-rest-api/internal/user/db"
	"go-rest-api/pkg/client/mongodb"
	"go-rest-api/pkg/logging"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {

	logger := logging.GetLogger()
	logger.Info("Create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	cfgMongo := cfg.MongoDB
	logger.Info("Create storage")
	client, err := mongodb.NewClient(context.Background(), cfgMongo.Host,
		cfgMongo.Port, cfgMongo.Username, cfgMongo.Password, cfgMongo.DB, cfgMongo.AuthDB)
	if err != nil {
		panic(err)
	}

	logger.Info("Initialization user storage")
	storage := db.NewStorage(client, cfgMongo.Collection, logger)

	logger.Info("Create user service")
	service := user.NewService(storage, logger)

	logger.Info("Create handler")
	handler := user.NewHandler(logger, service)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start app")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("Detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := filepath.Join(appDir, "app.sock")

		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("start listening unix socket %s", socketPath)

	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Infof("start listening on %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	logger.Fatal(server.Serve(listener))
}
