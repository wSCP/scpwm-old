package main

import (
	"fmt"
	"net"
	"os"

	"github.com/thrisp/scpwm/euclid/manager"
	"github.com/thrisp/scpwm/flarg"
	"github.com/thrisp/scpwm/version"
)

var (
	provided    *args
	pkgVersion  version.Version
	packageName string = "SCPWM euclid"
	versionTag  string = "No version tag supplied with compilation"
	versionHash string
	versionDate string
)

type args struct {
	SocketEnv     string `arg:"--socketEnv,help:Provides an env variable to locate the SCPWM socket. Euclid will attempt to read from 'SCPWM_SOCKET' by default."`
	socketPathTpl string `arg:"-"`
	SocketPath    string `arg:"--socketPath,help:Specifies a socket path directly. Default is '/tmp/scpwm_0_0-socket'"`
	ConfigHomeEnv string `arg:"--configEnv,help:Specifies an env variable for the config file home. The default is 'XDG_CONFIG_HOME'"`
	ConfigFile    string `arg:"-c,--configFile,help:Reads the main configuration from the given file of the configuration home. The default is 'euclid/euclidrc'"`
	Verbose       bool   `arg:"-v,--verbose,help:Allows for euclid to provided more and more detailed logging messages."`
	Version       bool   `arg:"--version,help:Prints the compiled program version and exit."`
}

func defaultArgs() *args {
	return &args{
		"SCPWM_SOCKET",
		"/tmp/scpwm%s_%d_%d-socket",
		"",
		"XDG_CONFIG_HOME",
		"euclid/euclidrc",
		false,
		false,
	}
}

func (a *args) socketPath() string {
	if a.SocketPath != "" {
		return a.SocketPath
	}
	var socketPath = os.Getenv(a.SocketEnv)
	if socketPath == "" {
		socketPath = fmt.Sprintf(a.socketPathTpl, "", 0, 0)
	}
	return socketPath
}

func (a *args) configPath() string {
	var pth string
	ch := os.Getenv(a.ConfigHomeEnv)
	if ch != "" {
		pth = fmt.Sprintf("%s/%s", ch, a.ConfigFile)
	} else {
		pth = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), ".config", a.ConfigFile)
	}
	return pth
}

func init() {
	provided = defaultArgs()
	flarg.MustParse(provided)
	pkgVersion = version.New(packageName, versionTag, versionHash, versionDate)
}

func main() {
	if provided.Version {
		fmt.Printf(pkgVersion.Fmt())
		os.Exit(0)
	}

	e := manager.New()

	sp := provided.socketPath()
	sckt, err := net.ListenUnix("unix", &net.UnixAddr{sp, "unix"})
	if err != nil {
		panic(err)
	}
	defer os.Remove(sp)
	defer e.Conn().Close()

	l := e.Looping(sckt)

	cp := provided.configPath()
	err = e.LoadConfig(cp)
	if err != nil {
		e.Println(err.Error())
	}

	e.Add("ConfigPath", cp)

	e.Println(pkgVersion.Fmt())

EVENT:
	for {
		select {
		case <-l.Pre:
			<-l.Post
		case msg := <-l.Comm:
			if provided.Verbose {
				e.Println(msg)
			}
		case sig := <-l.Sys:
			e.SignalHandler(sig)
		case <-l.Quit:
			break EVENT
		}
	}
}
