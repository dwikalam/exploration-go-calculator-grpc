package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/dwikalam/calcgorpc/internal/app/config"
	"github.com/dwikalam/calcgorpc/internal/app/pb"
	"github.com/dwikalam/calcgorpc/internal/app/server"
	"google.golang.org/grpc"
)

func Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	const (
		keepAliveTimeout time.Duration = time.Second * 60
	)

	var (
		cfg config.Config

		tcpConfig net.ListenConfig
		listener  net.Listener

		srv        *server.Server
		grpcServer *grpc.Server

		wg sync.WaitGroup

		err error
	)

	cfg, err = config.New()
	if err != nil {
		return fmt.Errorf("creating new config failed: %v", err)
	}

	// Setup Listener

	tcpConfig = net.ListenConfig{
		KeepAlive: keepAliveTimeout,
	}
	listener, err = tcpConfig.Listen(
		ctx,
		cfg.GetServerNetwork(),
		cfg.GetServerAddress(),
	)
	if err != nil {
		return fmt.Errorf(
			"%s listening failed on %s: %v",
			cfg.GetServerNetwork(),
			cfg.GetServerAddress(),
			err,
		)
	}

	fmt.Printf("%s listening on %s\n", cfg.GetServerNetwork(), cfg.GetServerAddress())

	// Setup Server

	go func() {
		grpcServer = grpc.NewServer()
		srv = server.New()

		pb.RegisterCalculatorServer(grpcServer, srv)

		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("failed serving on the listener: %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		if err = listener.Close(); err != nil {
			log.Printf("failed closing listener: %v\n", err)
		}

		fmt.Printf("listener closed\n")
	}()

	wg.Wait()

	return nil
}
