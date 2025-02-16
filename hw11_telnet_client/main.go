package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func printUsage() {
	fmt.Println("go-telnet\nUsage:")
	fmt.Println("  go-telnet --timeout=10s <host> <port>")
	fmt.Println("  go-telnet mysite.ru 8080")
	fmt.Println("  go-telnet --timeout=3s 1.1.1.1 123")
	fmt.Println("Named arguments:")
	flag.PrintDefaults()
}

func main() {
	timeout := flag.Duration("timeout", time.Duration(10*float64(time.Second)), "Connection establish timeout")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		printUsage()
	}
	address := net.JoinHostPort(args[0], args[1])
	logger := log.New(os.Stderr, "", 0)

	process(address, timeout, os.Stdin, os.Stdout, logger)
}

func process(address string, timeout *time.Duration, in io.ReadCloser, out io.Writer, logger *log.Logger) {
	client := NewTelnetClient(address, *timeout, in, out)

	err := client.Connect()
	if err != nil {
		os.Exit(1)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	eof := make(chan struct{}, 1)
	connectionClosed := make(chan struct{}, 1)
	go func() {
		if err := client.Send(); err != nil {
			logger.Printf("...Send error: %v", err)
		}
		eof <- struct{}{}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			logger.Printf("...Receive error: %v", err)
		}
		connectionClosed <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		logger.Println("...Ctrl+C")
		os.Exit(2)
	case <-eof:
	case <-connectionClosed:
		client.Close()
	}
}
