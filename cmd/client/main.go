package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/dwikalam/calcgorpc/internal/app/config"
	"github.com/dwikalam/calcgorpc/internal/app/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	if err := app(ctx); err != nil {
		log.Printf("running app failed: %v", err)

		os.Exit(1)
	}

}

func app(ctx context.Context) error {
	const (
		timeout time.Duration = time.Millisecond * 500
	)

	var (
		cancel context.CancelFunc

		cfg config.Config

		conn *grpc.ClientConn
		cc   pb.CalculatorClient

		err error
	)

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	cfg, err = config.New()
	if err != nil {
		return fmt.Errorf("creating new config failed: %v", err)
	}

	conn, err = grpc.NewClient(
		net.JoinHostPort(cfg.GetServerHost(), strconv.Itoa(cfg.GetServerPort())),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("creating new grpc client failed: %v", err)
	}

	cc = pb.NewCalculatorClient(conn)

	resp, err := cc.Add(ctx, &pb.AddRequest{A: 5, B: 10})
	if err != nil {
		log.Printf("failed to do Add RPC: %v\n", err)
	}

	fmt.Printf("Add RPC result: %d\n", resp.GetResult())

	return nil
}
