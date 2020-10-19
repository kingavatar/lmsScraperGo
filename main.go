package main

import(
    "fmt"
    "net"
    "log"
    "os"
    "os/signal"
    "syscall"
    "strings"
    "github.com/kingavatar/lmsScraperGo/scraper"
    "github.com/takama/daemon"
)

const (

		// name of the service
		name        = "kingScraperService"
		description = "kingavatar lmsScraper Service"

		// port which daemon should be listen
		port = ":9977"
  )
  
  var stdlog, errlog *log.Logger

	// Service has embedded daemon
	type Service struct {
		daemon.Daemon
  }
  
  // Manage by daemon commands or run the daemon
	func (service *Service) Manage() (string, error) {

		usage := "Usage: myservice install | remove | start | stop | status"

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
				return service.Stop()
			case "status":
				return service.Status()
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
				go handleConnection(conn,listener)
			case killSignal := <-interrupt:
				stdlog.Println("Got signal:", killSignal)
				stdlog.Println("Stoping listening on ", listener.Addr())
				listener.Close()
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
func handleConnection(c net.Conn,listener net.Listener){
        buf := make([]byte, 1512)
        nr, err := c.Read(buf)
        if err != nil {
            return
        }
        defer c.Close()
        data := buf[0:nr]
        fmt.Printf("Received: %v", string(data))
        var args string
        args=strings.TrimSpace(string(data))
        if(args=="events"){
          data =scraper.GetEvents()
        } else {
        data=scraper.GetUnknownResponse()
        }
        _, err = c.Write(data)
        if err != nil {
            log.Fatalln("Write Connection: " , err)
        }

}

//Init to start some stuff taken from daemon example code
func Init() {
    stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
    errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

//Start Function to run service Starts
func StartJob() error{
  Init()
  if err:= scraper.Start(); err != nil {
      return err
  }
  // service.Manage()
  // fmt.Println("The Service has Started")
  return nil
}

//Stop Function to run service Stops
func Stop() error{
  if err:= scraper.Stop(); err != nil {
      return err
  }
  // fmt.Println("The Service has Started")
  return nil
}

func main(){
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
    srv, err := daemon.New(name, description, daemon.SystemDaemon) //dependencies...)
    if err != nil {
        errlog.Println("Error: ", err)
        os.Exit(1)
    }
    service := &Service{srv}
    StartJob()
    status, err := service.Manage()
    if err != nil {
        errlog.Println(status, "\nError: ", err)
        os.Exit(1)
    }
    fmt.Println(status)
}