package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

var (
	Logger          = log.New(os.Stderr, "KETER: ", log.Ldate|log.Lmicroseconds)
	CONFIG_HOME_ENV = "XDG_CONFIG_HOME"
	CONFIG_PATH     = "keter/keterrc"
	ChainExpiry     int
	ConfigPath      string
	verbose         bool
)

func defaultConfigPath() string {
	var pth string
	configHome := os.Getenv(CONFIG_HOME_ENV)
	if configHome != "" {
		pth = fmt.Sprintf("%s/%s", configHome, CONFIG_PATH)
	} else {
		pth = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), ".config", CONFIG_PATH)
	}
	return pth
}

func SignalHandler(h Handlr, s os.Signal) {
	msg := new(bytes.Buffer)
	switch s {
	case syscall.SIGINT:
		Logger.Println("SIGINT")
		os.Exit(0)
	case syscall.SIGHUP:
		msg.WriteString("Got signal SIGHUP, reconfiguring....\n")
		chains, err := LoadConfig(ConfigPath)
		if err != nil {
			msg.WriteString(fmt.Sprintf("error while loading config: %s\n", err.Error()))
		}
		err = Configure(h, chains)
		if err != nil {
			msg.WriteString(fmt.Sprintf("error while configuring: %s\n", err))
		}
	default:
		Logger.Println(fmt.Sprintf("received %v", s))

	}
	if verbose && msg.Len() != 0 {
		Logger.Println(msg.String())
	}
}

func main() {
	chains, err := LoadConfig(ConfigPath)
	if err != nil {
		Logger.Fatalf("configuration loading error: %s", err.Error())
	}

	hndl, err := NewHandlr("")
	if err != nil {
		Logger.Fatalf("handler configuration error: %s", err.Error())
	}

	err = Configure(hndl, chains)
	if err != nil {
		Logger.Fatalf("key chain configuration error: %s", err.Error())
	}

	before, after, quit, signals := Loop(hndl)

EVENTLOOP:
	for {
		select {
		case <-before:
			<-after
		case sig := <-signals:
			SignalHandler(hndl, sig)
		case <-quit:
			break EVENTLOOP
		}
	}
}

func init() {
	flag.IntVar(&ChainExpiry, "timeout", 2, "Timeout in seconds for the recording of chord chains.")
	flag.StringVar(&ConfigPath, "config", defaultConfigPath(), "Read the main configuration from the given file.")
	flag.BoolVar(&verbose, "verbose", verbose, "Verbose logging messages, default is false.")
	flag.Parse()
}
