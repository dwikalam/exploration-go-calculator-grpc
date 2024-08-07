package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
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
		wg     sync.WaitGroup
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
	defer conn.Close()

	cc = pb.NewCalculatorClient(conn)

	wg.Add(1)
	go func() {
		var (
			resp *pb.CalculationResponse

			err error
		)

		defer wg.Done()

		resp, err = cc.Add(ctx, &pb.CalculationRequest{A: 5, B: 10})
		if err != nil {
			log.Printf("failed doing Add RPC: %v\n", err)

			return
		}

		fmt.Printf("Add RPC result: %.2f\n", resp.GetResult())
	}()

	wg.Add(1)
	go func() {
		var (
			resp *pb.CalculationResponse

			err error
		)

		defer wg.Done()

		resp, err = cc.Substract(ctx, &pb.CalculationRequest{A: 5, B: 10})
		if err != nil {
			log.Printf("failed doing Substract RPC: %v\n", err)

			return
		}

		fmt.Printf("Substract RPC result: %.2f\n", resp.GetResult())
	}()

	wg.Add(1)
	go func() {
		var (
			resp *pb.CalculationResponse

			err error
		)

		defer wg.Done()

		resp, err = cc.Multiply(ctx, &pb.CalculationRequest{A: 5, B: 10})
		if err != nil {
			log.Printf("failed doing Multiply RPC: %v\n", err)

			return
		}

		fmt.Printf("Multiply RPC result: %.2f\n", resp.GetResult())
	}()

	wg.Add(1)
	go func() {
		var (
			resp *pb.CalculationResponse

			err error
		)

		defer wg.Done()

		resp, err = cc.Divide(ctx, &pb.CalculationRequest{A: 5, B: 0})
		if err != nil {
			log.Printf("failed doing Divide RPC: %v\n", err)

			return
		}

		fmt.Printf("Divide RPC result: %.2f\n", resp.GetResult())
	}()

	wg.Add(1)
	go func() {
		var (
			resp *pb.CalculationResponse

			err error
		)

		defer wg.Done()

		resp, err = cc.Divide(ctx, &pb.CalculationRequest{A: 5, B: 10})
		if err != nil {
			log.Printf("failed doing Divide RPC: %v\n", err)

			return
		}

		fmt.Printf("Divide RPC result: %.2f\n", resp.GetResult())
	}()

	wg.Add(1)
	go func() {
		var (
			resp *pb.CalculationResponse

			err error
		)

		defer wg.Done()

		resp, err = cc.Sum(ctx, &pb.NumbersRequest{Numbers: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}})
		if err != nil {
			log.Printf("failed doing Sum RPC: %v\n", err)

			return
		}

		fmt.Printf("Sum RPC result: %.2f\n", resp.GetResult())
	}()

	wg.Wait()

	return nil
}
