# lmsScraperGo

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kingavatar/lmsScraperGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/kingavatar/lmsScraperGo)](https://goreportcard.com/report/github.com/kingavatar/lmsScraperGo)
[![GoDoc](https://godoc.org/github.com/kingavatar/lmsScraperGo/scraper?status.svg)](https://godoc.org/github.com/kingavatar/lmsScraperGo/scraper)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/kingavatar/lmsScraperGo/scraper)](https://pkg.go.dev/github.com/kingavatar/lmsScraperGo/scraper)

[![forthebadge made-with-go](https://forthebadge.com/images/badges/made-with-go.svg)](https://golang.org/)

A Daemon Scraping Service for `LMS`(Learning Management System) written for my college IIITB but can serve as template for other `LMS` portals as well.

This Application is written in `Go` (`golang`) and scrapes for assignments and announcements of favorite(`starred` in LMS) courses.

## Features

-   same cookie and session when service starts.
-   can be run as system daemon process or temporary application process.
-   output through socket call(default `port:9977`) so can be read on other devices on the network as well.
-   color support for `conky` and terminal( with true color).
-   utilizes inbuilt `Go` concurrency.

## Installation

Make sure to have latest `go` version installed.


```zsh
$ go version
go version go1.15.2 linux/amd64
```

To Fetch the Source.

To install in a local directory temporarily change `GOPATH` variable.

```
$ export GOPATH=$HOME/{path} ; go get github.com/kingavatar/lmsScraperGo

```
Before Building please change the login authentication Variables in `login.go` file present in `scraper` folder.

```go
var (
	username = "XXXXXX"
	password = "XXXXXX"
)
```
Then Build the source and copy the binary to local `bin` folder[optional].

```bash
$ go build github.com/kingavatar/lmsScraperGo
$ cp lmsScraperGo ~/.local/bin/
```

First you have to install the service.

```zsh
$ sudo lmsScraperGo install
Install kingavatar lmsScraper Service:                                  [  OK  ]
```

Root is required to install the daemon service in `systemd/system` folder. The service name is `kingScraper`.


Then start the service.

```zsh
$ sudo lmsScraperGo start
Starting kingavatar lmsScraper Service:                                 [  OK  ]
```

Then to see the service pid run with the status command.
```zsh
$ sudo lmsScraperGo status
Service (pid  pidno) is running...
```

To see the actual log outputs see using systemctl or systemd command.
```zsh
$ sudo systemctl status kingScraper
```

To stop the service
```zsh
$ sudo lmsScraperGo stop
Stopping kingavatar lmsScraper Service:                                 [  OK  ]
```

To remove the service from the system
```zsh
$ sudo lmsScraperGo remove
Removing kingavatar lmsScraper Service:                                 [  OK  ]
```

You can also run it temporarily
```zsh
$ lmsScraperGo
```
Please stop the daemon service to avoiding binding to same port or you will receive the following error.

```zsh
Error:  listen tcp :9977: bind: address already in use
``` 

You can also run it in debug mode to get more logs.
```zsh
$ lmsScraperGo -d
```
## Usage

You can use netcat or other services which connects and reads socket to get the scraped output.

To get `events`(Assignments) and `announcements` for `starred` Courses.
```zsh
$ echo "events" | nc localhost 9977
$ echo "announcements" | nc localhost 9977
```
To get Terminal( with true color) colorful output.

```zsh
$ echo -e $(echo "eventsterm" | nc localhost 9977)
$ echo -e $(echo "announcementsterm" | nc localhost 9977)
```

For Conky you can use the below example for getting output every 5 minutes.
```conky
execpi 300 echo "eventsconky" | nc localhost 9977
execpi 300 echo "announcementsconky" | nc localhost 9977
```
