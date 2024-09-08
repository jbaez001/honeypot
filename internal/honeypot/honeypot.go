package honeypot

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/containrrr/shoutrrr"
	"github.com/jbaez001/honeypot/internal/config"
)

var honeypotConfig config.HoneypotConfig

func Shout(urls []string, msg string) {
	sender, err := shoutrrr.CreateSender(urls...)
	if err != nil {
		log.Println("Error creating sender: ", err)
		return
	}

	sender.Send(msg, nil)

}

func Start(cfg config.HoneypotConfig) {
	honeypotConfig = cfg

	var wg sync.WaitGroup

	// start the honeypots (fake servers)
	for _, pot := range honeypotConfig.Honeypots {
		wg.Add(1)
		go func(s config.Honeypot) {
			defer wg.Done()
			handleConnection(pot)
		}(pot)
	}

	// block until infinity 'n beyond
	wg.Wait()
	log.Println("Honeypot stopped")
}

func handleConnection(cfg config.Honeypot) {
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Printf("[%s] ERROR starting on port %s: %v\n", cfg.Protocol, cfg.Port, err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Printf("[%s] ERROR closing listener: %v\n", cfg.Protocol, err)
		}
	}(listener)

	log.Printf("[%s] STARTED on %s\n", cfg.Protocol, cfg.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[%s] ERROR while trying to accept connection from: %v\n", cfg.Protocol, err)
			continue
		}
		msg := fmt.Sprintf("<%s> [%s] <- %s", honeypotConfig.Name, cfg.Protocol, conn.RemoteAddr())
		err = conn.Close()
		if err != nil {
			log.Printf("[%s] ERROR closing connection: %v\n", cfg.Protocol, err)
		}
		log.Println(msg)

		if len(honeypotConfig.ShoutUrls) > 0 {
			Shout(honeypotConfig.ShoutUrls, msg)
		}

		// if the server is fragile, stop listening
		if cfg.Fragile {
			log.Printf("[%s] FRAGILE server, stopping...\n", cfg.Protocol)
			return
		}
	}
}
