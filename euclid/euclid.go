package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/thrisp/scpwm/euclid/manager"
)

var (
	socketEnv     string = "SCPWM_SOCKET"
	socketPathTpl string = "/tmp/scpwm%s_%d_%d-socket"
	socketPath    string
	ConfigHomeEnv string = "XDG_CONFIG_HOME"
	ConfigFile    string = "euclid/euclidrc"
	ConfigPath    string
	verbose       bool
)

func defaultSocketPath() string {
	var socketPath = os.Getenv(socketEnv)
	if socketPath == "" {
		socketPath = fmt.Sprintf(socketPathTpl, "", 0, 0)
	}
	return socketPath
}

func main() {
	e := manager.New()
	sckt, err := net.ListenUnix("unix", &net.UnixAddr{socketPath, "unix"})
	if err != nil {
		panic(err)
	}
	defer os.Remove(socketPath)

	defer e.Conn().Close()

	l := e.Looping(sckt)

	err = e.LoadConfig(ConfigPath)

	if err != nil {
		e.Println(err.Error())
	}

EVT:
	for {
		select {
		case <-l.Pre:
			<-l.Post
		case msg := <-l.Comm:
			if verbose {
				e.Println(msg)
			}
		case <-l.Quit:
			break EVT
		}
	}

	e.Println("EXITING......\n")
}

func defaultConfigPath() string {
	var pth string
	ch := os.Getenv(ConfigHomeEnv)
	if ch != "" {
		pth = fmt.Sprintf("%s/%s", ch, ConfigFile)
	} else {
		pth = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), ".config", ConfigFile)
	}
	return pth
}

func init() {
	flag.StringVar(&ConfigPath, "config", defaultConfigPath(), "Reads the main configuration from the given file. The default location is 'XDG_CONFIG_HOME/euclid/euclidrc'")
	flag.StringVar(&socketEnv, "e", socketEnv, "Reads the socket from the given env variable, scpwm will attempt to read from 'SCPWM_SOCKET' by default")
	flag.StringVar(&socketPath, "p", defaultSocketPath(), "Reads the socket from the given path. Default is '/tmp/scpwm_0_0-socket'")
	flag.BoolVar(&verbose, "v", verbose, "Verbose logging messages, default is false.")
	flag.Parse()
}
