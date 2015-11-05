package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/thrisp/scpwm/flarg"
	"github.com/thrisp/scpwm/version"
)

var (
	provided    *args
	Logger      = log.New(os.Stderr, "KETER: ", log.Ldate|log.Lmicroseconds)
	pkgVersion  version.Version
	packageName string = "SCPWM keter"
	versionTag  string = "No version tag supplied with compilation"
	versionHash string
	versionDate string
)

type args struct {
	ConfigEnv   string        `arg:"--configEnv,help:Provide an env variable to locate the home directory of the configuration file. Default is 'XDG_CONFIG_HOME.'"`
	ConfigPath  string        `arg:"--configPath,help:Provide the path of the config file in the config home directory. Default is 'keter/keterrc'"`
	ChainExpiry time.Duration `arg:"-t,--timeout,help:Timeout in seconds for the recording of chord chains. Default is 2."`
	Verbose     bool          `arg:"--verbose,help:Verbose logging messages."`
	Version     bool          `arg:"-v,--version,help:Print the compiled program version and exit."`
}

func defaultArgs() *args {
	return &args{
		"XDG_CONFIG_HOME",
		"keter/keterrc",
		2 * time.Second,
		false,
		false,
	}
}

func (a *args) configPath() string {
	var pth string
	configHome := os.Getenv(a.ConfigEnv)
	if configHome != "" {
		pth = fmt.Sprintf("%s/%s", configHome, a.ConfigPath)
	} else {
		pth = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), ".config", a.ConfigPath)
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
		chains, err := LoadConfig(provided.configPath())
		if err != nil {
			msg.WriteString(fmt.Sprintf("error while loading config: %s\n", err.Error()))
		}
		err = Configure(h, chains)
		if err != nil {
			msg.WriteString(fmt.Sprintf("error while configuring: %s\n", err))
		}
	default:
		Logger.Println(fmt.Sprintf("received signal %v", s))
	}
	if provided.Verbose && msg.Len() != 0 {
		Logger.Println(msg.String())
	}
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
	chains, err := LoadConfig(provided.configPath())
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
