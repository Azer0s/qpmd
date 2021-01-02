package internals

import (
	"fmt"
	"github.com/Azer0s/qpmd"
	"github.com/takama/daemon"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	name        = "qpmd"
	description = "quacktor port mapping daemon"
	port        = ":7161"
	timeout     = time.Second * 30
)

var stdLog, errLog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {
	usage := "Usage: qpmd install | remove | start | stop | status | version"

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		case "version":
			return fmt.Sprintln(description) +
				fmt.Sprint("Version:  ", qpmd.VERSION), nil
		default:
			return usage, nil
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	stdLog.Println("Starting listening on port 7161")

	listener, err := net.Listen("tcp", port)

	if err != nil {
		return "Couldn't start listening on port 7161", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			stdLog.Printf("Received %s", killSignal.String())
			stdLog.Printf("Closing listener on 7161")
			err = listener.Close()

			if err != nil {
				return "There was a problem while closing the listener", err
			}

			stdLog.Printf("Closed listener on 7161")

			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}

			return "Daemon was killed", nil
		}
	}
}

// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func init() {
	stdLog = log.New(os.Stdout, "[qpmd] ", log.Ldate|log.Ltime)
	errLog = log.New(os.Stderr, "[qpmd] ", log.Ldate|log.Ltime)
}

func Launch() {
	var srv daemon.Daemon
	var err error

	switch runtime.GOOS {
	case "darwin":
		srv, err = daemon.New(name, description, daemon.UserAgent)
	case "freebsd":
		srv, err = daemon.New(name, description, daemon.SystemDaemon)
	case "linux":
		srv, err = daemon.New(name, description, daemon.SystemDaemon)
	case "windows":
		srv, err = daemon.New(name, description, daemon.SystemDaemon)
	}

	if err != nil {
		errLog.Println("Error: ", err)
		os.Exit(1)
	}

	service := &Service{srv}
	status, err := service.Manage()

	if err != nil {
		errLog.Println(status, "\nError: ", err)
		os.Exit(1)
	}

	fmt.Println(status)
}
