package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/deeper-x/tcpstream/errorcs"
	"github.com/deeper-x/tcpstream/models"
	"github.com/deeper-x/tcpstream/settings"
	"github.com/deeper-x/tcpstream/utils"
)

var (
	clientsMu sync.Mutex
	clients   = make(map[net.Conn]*models.Client)
)

func Run() error {
	err := os.Remove(settings.SocketPath)
	if err != nil {
		return err
	}

	l, err := net.Listen("unix", settings.SocketPath)
	if err != nil {
		return err
	}
	defer l.Close()

	log.Println("chatgo listening on", settings.SocketPath)

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	if !scanner.Scan() {
		os.Exit(errorcs.SCANNING)
	}

	welcomeLine := scanner.Text()
	ok := utils.IsAllowed(welcomeLine)

	if !ok {
		os.Exit(errorcs.DENIED)
	}

	client := &models.Client{
		Name: welcomeLine[settings.LenMOTD:],
		Conn: conn,
	}

	clientsMu.Lock()
	clients[conn] = client
	clientsMu.Unlock()

	err := broadcast(client, fmt.Sprintf("%s joined\n", client.Name))
	if err != nil {
		os.Exit(errorcs.JOIN)
	}

	for scanner.Scan() {
		msg := scanner.Text()
		err := broadcast(client, fmt.Sprintf("%s says: %s\n", client.Name, msg))
		if err != nil {
			os.Exit(errorcs.MSG)
		}

	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()

	err = broadcast(client, fmt.Sprintf("%s left...\n", client.Name))
	if err != nil {
		os.Exit(errorcs.LEFT)
	}
}

func broadcast(sender *models.Client, msg string) error {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for _, c := range clients {
		if c != sender {
			_, err := c.Conn.Write([]byte(msg))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
