package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"tailscale.com/ipn"
	"tailscale.com/tsnet"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s host port\n", os.Args[0])
		os.Exit(2)
	}
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	host := os.Args[1]
	port := os.Args[2]
	addr := host + ":" + port

	s := &tsnet.Server{
		Hostname:  "ssh-proxycommand",
		Dir:       os.Getenv("TSNET_STATE_DIR"),
		Ephemeral: true,
	}

	if s.Dir == "" {
		s.Dir = os.ExpandEnv("$HOME/.local/state/tsnet-ssh-proxy")
	}

	if _, err := s.Up(context.Background()); err != nil {
		return fmt.Errorf("tsnet up: %w", err)
	}

	lc, err := s.LocalClient()
	if err != nil {
		log.Printf("failed to get local client: %v", err)
	}

	if lc != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		watcher, err := lc.WatchIPNBus(ctx, ipn.NotifyInitialNetMap)
		if err == nil {
			for {
				status, err := lc.Status(ctx)
				if err != nil {
					log.Printf("failed to get status: %v", err)
					break
				}
				if status.Self != nil && status.Self.Online {
					break
				}

				_, err = watcher.Next()
				if err != nil {
					break
				}
			}
			watcher.Close()
		} else {
			log.Printf("failed to watch IPN bus: %v, proceeding to dial", err)
		}
		cancel()
	}

	defer func() {
		log.Printf("cleanup: starting...")
		log.Printf("cleanup: closing tsnet server...")
		if err := s.Close(); err != nil {
			log.Printf("cleanup: failed to close server: %v", err)
		}
		log.Printf("cleanup: finished")
	}()

	conn, err := s.Dial(context.Background(), "tcp", addr)
	if err != nil {
		return fmt.Errorf("tsnet dial %s: %w", addr, err)
	}
	defer conn.Close()

	errc := make(chan error, 2)

	go func() {
		_, err := io.Copy(conn, os.Stdin)
		errc <- err
	}()

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		errc <- err
	}()

	return <-errc
}
