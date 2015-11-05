package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/thrisp/scpwm/flarg"
	"github.com/thrisp/scpwm/version"
)

var (
	provided    *args
	logger      = log.New(os.Stderr, "[SCPC] ", log.Ldate|log.Lmicroseconds)
	pkgVersion  version.Version
	packageName string = "SCPWM scpc"
	versionTag  string = "No version tag supplied with compilation"
	versionHash string
	versionDate string
)

type args struct {
	SocketEnv     string        `arg:"-s,--socket"`
	socketPathTpl string        `arg:"-"`
	SocketPath    string        `arg:"-p,--path"`
	Timeout       time.Duration `arg:"-t,--timeout"`
	Verbose       bool          `arg:"--verbose"`
	Version       bool          `arg:"-v,--version"`
	Command       []string      `arg:"positional"`
}

func defaultArgs() *args {
	return &args{
		"SCPWM_SOCKET",
		"/tmp/scpwm%s_%d_%d-socket",
		"",
		(2 * time.Second),
		false,
		false,
		nil,
	}
}

func (a *args) socketPath() string {
	var socketPath = os.Getenv(a.SocketEnv)
	if socketPath == "" {
		socketPath = fmt.Sprintf(a.socketPathTpl, "", 0, 0)
	}
	return socketPath
}

func removable(index int, remove []int) bool {
	for _, r := range remove {
		if index == r {
			return true
		}
	}
	return false
}

func send(command []string) []byte {
	stop := len(command) - 1

	b := new(bytes.Buffer)
	for i, v := range command {
		b.WriteString(v)
		if i != stop {
			b.WriteRune(' ')
		}
	}
	return b.Bytes()
}

func responded(r Response, sent []byte, verbose bool) Response {
	if verbose {
		logger.Println(fmt.Sprintf("euclid returned: `%s` to sent message: `%s`", r, sent))
	}
	return r
}

func timedout(to time.Duration, verbose bool) Response {
	if verbose {
		logger.Println("timeout, no response from euclid after %d seconds", to)
	}
	return timeout
}

type Response int

const (
	unknown Response = iota
	failure
	syntax
	length
	timeout
	success
)

func (r Response) String() string {
	switch r {
	case unknown:
		return "unknown"
	case failure:
		return "failure"
	case syntax:
		return "syntax"
	case length:
		return "length"
	case timeout:
		return "timeout"
	case success:
		return "success"
	}
	return ""
}

var stringResponse = map[string]Response{
	"unknown": unknown,
	"failure": failure,
	"syntax":  syntax,
	"length":  length,
	"timeout": timeout,
	"success": success,
}

func response(c io.Reader, provided *args, sent []byte) Response {
	buf := make([]byte, 1024)
	t := time.After(provided.Timeout)
	for {
		select {
		case <-t:
			return timedout(provided.Timeout, provided.Verbose)
		default:
			n, err := c.Read(buf[:])
			if err != nil {
				logger.Println(err)
				return failure
			}
			if resp, ok := stringResponse[string(buf[0:n])]; ok {
				return responded(resp, sent, provided.Verbose)
			}
		}
	}
	return timeout
}

func exit(c *net.UnixConn, code int) {
	c.Close()
	os.Remove("/tmp/scpc")
	os.Exit(code)
}

func main() {
	if provided.Version {
		fmt.Printf(pkgVersion.Fmt())
		os.Exit(0)
	}

	t := "unix"
	laddr := net.UnixAddr{"/tmp/scpc", t}
	conn, err := net.DialUnix(t, &laddr, &net.UnixAddr{provided.socketPath(), t})
	if err != nil {
		panic(err)
	}

	sent := send(provided.Command)

	_, err = conn.Write(sent)
	if err != nil {
		panic(err)
	}

	rsp := response(conn, provided, sent)
	switch rsp {
	case success:
		exit(conn, 0)
	case failure, length, timeout:
		exit(conn, 1)
	case syntax:
		exit(conn, 2)
	default:
		exit(conn, 3)
	}
}

func init() {
	provided = defaultArgs()
	flarg.MustParse(provided)
	pkgVersion = version.New(packageName, versionTag, versionHash, versionDate)
}
