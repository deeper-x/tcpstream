package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/deeper-x/chatgo/errorcs"
	"github.com/deeper-x/chatgo/settings"
	"github.com/deeper-x/chatgo/utils"
)

func Run() error {
	if len(os.Args) < 2 {
		return errors.New(settings.UsageStr)
	}

	username := utils.GetUsername()

	conn, err := net.Dial("unix", settings.SocketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s%s\n", settings.MOTD, username)

	// Reader from socket
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		os.Exit(errorcs.SOCKET)
	}()

	// Reader from stdin
	stdin := bufio.NewScanner(os.Stdin)

	for stdin.Scan() {
		fmt.Fprintln(conn, stdin.Text())
	}

	return nil
}
