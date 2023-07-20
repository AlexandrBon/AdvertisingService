package main

import (
	"advertisingService/internal/adApp"
	"advertisingService/internal/adapters/adrepo"
	"advertisingService/internal/adapters/userrepo"
	grpcPort "advertisingService/internal/ports/grpc"
	"advertisingService/internal/ports/httpgin"
	"advertisingService/internal/userApp"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	adRepository := adrepo.New()
	userRepository := userrepo.New()
	ginServer := httpgin.NewHTTPServer(":18080", adApp.NewApp(adRepository, userRepository), userApp.NewApp(userRepository))
	err := ginServer.Listen()
	if err != nil {
		panic(err)
	}

	port := ":18081"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	eg, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	grpcServer := grpcPort.NewGRPCServer(lis, adApp.NewApp(adRepository, userRepository), userApp.NewApp(userRepository))

	eg.Go(func() error {
		log.Printf("starting main server, listening on %s\n", port)
		defer log.Printf("close main server listening on %s\n", port)

		errCh := make(chan error)

		defer func() {
			grpcServer.GetServer().GracefulStop()
			_ = lis.Close()

			close(errCh)
		}()

		go func() {
			if err := grpcServer.Listen(); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("main server can't listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}

	log.Println("servers were successfully shutdown")
}
