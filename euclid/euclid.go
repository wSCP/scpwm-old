package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/thrisp/scpwm/euclid/manager"
	"github.com/thrisp/scpwm/version"
)

var (
	socketEnv     string = "SCPWM_SOCKET"
	socketPathTpl string = "/tmp/scpwm%s_%d_%d-socket"
	socketPath    string
	ConfigHomeEnv string = "XDG_CONFIG_HOME"
	ConfigFile    string = "euclid/euclidrc"
	ConfigPath    string
	verbose       bool
	pkgVersion    version.Version
	packageName   string = "SCPWM euclid"
	versionTag    string = "No version tag supplied with compilation"
	versionHash   string
	versionDate   string
	callVersion   bool
)

func defaultSocketPath() string {
	var socketPath = os.Getenv(socketEnv)
	if socketPath == "" {
		socketPath = fmt.Sprintf(socketPathTpl, "", 0, 0)
	}
	return socketPath
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
	flag.BoolVar(&verbose, "verbose", verbose, "Verbose logging messages, default is false.")
	flag.BoolVar(&callVersion, "version", callVersion, "Print the package version.")
	flag.Parse()
	pkgVersion = version.New(packageName, versionTag, versionHash, versionDate)
}

func main() {
	if callVersion {
		fmt.Printf(pkgVersion.Fmt())
		os.Exit(0)
	}
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

	e.Println("running...")
	e.Println(pkgVersion.Fmt())

EVENT:
	for {
		select {
		case <-l.Pre:
			<-l.Post
		case msg := <-l.Comm:
			if verbose {
				e.Println(msg)
			}
		case sig := <-l.Sys:
			e.SignalHandler(sig)
		case <-l.Quit:
			break EVENT
		}
	}
}
