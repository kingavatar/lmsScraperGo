package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kingavatar/lmsScraperGo/scraper"
	"github.com/takama/daemon"
)

const (

	// name of the service
	name        = "kingScraper"
	description = "kingavatar lmsScraper Service"

	// port which daemon should be listen
	port = ":9977"
)

var stdlog, errlog *log.Logger

//Debug to set debugging mode
var Debug *bool

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: lmsScrpaerGo install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			StartJob()
			return service.Start()
		case "stop":
			StopJob()
			return service.Stop()
		case "status":
			return service.Status()
		case "-d":
			break
		default:
			return usage, nil
		}
	}

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
			go handleConnection(conn, listener)
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			StopJob()
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
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
func handleConnection(c net.Conn, listener net.Listener) {
	buf := make([]byte, 1512)
	nr, err := c.Read(buf)
	if err != nil {
		return
	}
	defer c.Close()
	data := buf[0:nr]
	fmt.Printf("Received: %v", string(data))
	var args string
	args = strings.TrimSpace(string(data))

	switch args {
	case "events":
		data = scraper.GetEvents()
	case "announcements":
		data = scraper.GetAnnouncements()
	case "eventsterm":
		data = scraper.GetEventsTerm()
	case "eventsconky":
		data = scraper.GetEventsConky()
	case "announcementsterm":
		data = scraper.GetAnnouncementsTerm()
	case "announcementsconky":
		data = scraper.GetAnnouncementsConky()
	default:
		data = scraper.GetUnknownResponse()
	}
	_, err = c.Write(data)
	if err != nil {
		log.Fatalln("Write Connection: ", err)
	}

}

//StartJob Function to run service Starts
func StartJob() error {
	if err := scraper.Start(); err != nil {
		return err
	}
	// service.Manage()
	// fmt.Println("The Service has Started")
	return nil
}

//StopJob Function to run service Stops
func StopJob() error {
	if err := scraper.Stop(); err != nil {
		return err
	}
	return nil
}

func main() {
	//    if err:= scraper.Start(); err != nil {
	//     panic(fmt.Errorf("Something is Wrong"))
	//   }
	//   ln, err := net.Listen("tcp", ":9997")
	//   if err != nil {
	//     log.Fatalln("Net Connection Error")
	//   }
	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		// handle error
	// 	}
	// 	go handleConnection(conn)
	// }
	Debug = flag.Bool("d", false, "run in debug mode")
	flag.Parse()
	scraper.SetDebugMode(Debug)
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	srv, err := daemon.New(name, description, daemon.SystemDaemon) //dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	service.SetTemplate("[Unit]\nDescription={{.Description}}\nRequires={{.Dependencies}}\nAfter={{.Dependencies}}\n\n[Service]\nPIDFile=/var/run/{{.Name}}.pid\nExecStartPre=/bin/rm -f /var/run/{{.Name}}.pid\nExecStartPre=/bin/sh -c 'until ping -c1 archlinux.org; do sleep 1; done;'\nExecStart={{.Path}} {{.Args}}\nRestart=on-failure\n\n[Install]\nWantedBy=multi-user.target")
	if len(os.Args) < 2 || os.Args[1] == "-d" {
		err = StartJob()
		if err != nil {
			errlog.Println(err)
			os.Exit(1)
		}
	}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
