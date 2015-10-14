package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	CONFIG_HOME_ENV = "XDG_CONFIG_HOME"
	//KETER_SHELL_ENV = "KETER_SHELL"
	///SHELL_ENV       = "SHELL"
	CONFIG_PATH = "keter/keterrc"
	//socketPath      string
	//SOCKET_ENV      = "SCPWM_SOCKET"
	//SOCKET_PATH_TPL = "/tmp/scpwm%s_%d_%d-socket"
	ChainExpiry int
	ConfigPath  string
	Logger      = log.New(os.Stderr, "KETER: ", log.Lshortfile)
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

func main() {
	chains, err := LoadConfig(ConfigPath)

	if err != nil {
		Logger.Fatalf("configuration loading error: %s", err.Error())
	}

	X, err := Configure(chains)

	if err != nil {
		Logger.Fatalf("configuration error: %s", err.Error())
	}

	//socket path

	//fifo path

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP)
	b, a, q := Main(X)

EVENTLOOP:
	for {
		select {
		case <-b:
			<-a
		case sig := <-s:
			fmt.Printf("Got A HUP Signal %s! Now Reloading Conf....\n", sig)
		case <-q:
			break EVENTLOOP
		}
	}
}

func init() {
	flag.IntVar(&ChainExpiry, "timeout", 2, "Timeout in seconds for the recording of chord chains.")
	flag.StringVar(&ConfigPath, "config", defaultConfigPath(), "Read the main configuration from the given file.")
	flag.Parse()
}
