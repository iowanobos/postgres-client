package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/iowanobos/postgres-client/domain/account"
	"github.com/iowanobos/postgres-client/postgres"
	"github.com/sirupsen/logrus"
)

const DSN = "postgres://root:root@localhost:5432/test?sslmode=disable"

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	client, err := postgres.New(ctx, postgres.Options{ConnString: DSN})
	if err != nil {
		log.Fatal("Create postgres client failed. Error: ", err)
	}

	service := account.NewService(
		client.TransactionManager(),
		account.NewRepository(client),
	)

	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err = service.Test(ctx); err != nil {
				logrus.Error("Test service failed. Error: ", err)
			}
		}()
	}
	wg.Wait()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	select {
	case <-sigc:
		logrus.Info("Start shutdowning")
		cancel()
		client.Close()
	}
	logrus.Info("Application shut downing...")
}
