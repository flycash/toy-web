package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

func startServer(addr string, close chan struct{}) error {
	svc := http.Server{
		Addr: addr,
	}
	go func() {
		<-close
		fmt.Println("shutdown server on addr %s", addr)
		svc.Shutdown(context.Background())
	}()
	return svc.ListenAndServe()
}

func listenToSig(close chan struct{}) error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	select {
	case <-close:
		fmt.Println("stop waiting for system singanl")
	case <-sigint:
		fmt.Println("system int signal")
		return errors.New("systint")
	}
	return nil
}

func main() {
	gp, _ := errgroup.WithContext(context.Background())
	shutdownFlag := make(chan struct{})
	gp.Go(func() error {

		return startServer("127.0.0.1:8080", shutdownFlag)
	})
	gp.Go(func() error {
		return startServer("127.0.0.1:8081", shutdownFlag)
	})
	gp.Go(func() error {
		return listenToSig(shutdownFlag)
	})
	if err := gp.Wait(); err != nil {
		fmt.Println("error in group")
	}
	close(shutdownFlag)

	// to see println for debugging
	time.Sleep(1 * time.Second)
}
