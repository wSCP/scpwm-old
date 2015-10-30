package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	socketEnv       string = "SCPWM_SOCKET"
	socketPathTpl   string = "/tmp/scpwm%s_%d_%d-socket"
	socketPath      string
	responseTimeout time.Duration = 2
	verbose         bool

	logger = log.New(os.Stderr, "[SCPC] ", log.Ldate|log.Lmicroseconds)
)

func defaulSocketPath() string {
	var socketPath = os.Getenv(socketEnv)
	if socketPath == "" {
		socketPath = fmt.Sprintf(socketPathTpl, "", 0, 0)
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

func send(args []string) []byte {
	var r []int
	for i, v := range args {
		switch v {
		case "-verbose", "-timeout", "-path", "-env":
			r = append(r, i, i+1)
		}
	}

	var send []string

	for i, v := range args {
		if !removable(i, r) {
			send = append(send, v)
		}
	}

	idx := len(send) - 1

	b := new(bytes.Buffer)
	for i, v := range send {
		b.WriteString(v)
		if i != idx {
			b.WriteRune(' ')
		}
	}
	return b.Bytes()
}

func responded(r Response, sent []byte) Response {
	if verbose {
		logger.Println(fmt.Sprintf("euclid returned: `%s` to sent message: `%s`", r, sent))
	}
	return r
}

func timedout() Response {
	if verbose {
		logger.Println("timeout, no response from euclid after %d seconds", responseTimeout)
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

func response(c io.Reader, sent []byte) Response {
	buf := make([]byte, 1024)
	t := time.After(responseTimeout * time.Second)
	for {
		select {
		case <-t:
			return timedout()
		default:
			n, err := c.Read(buf[:])
			if err != nil {
				logger.Println(err)
				return failure
			}
			if resp, ok := stringResponse[string(buf[0:n])]; ok {
				return responded(resp, sent)
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
	t := "unix"
	laddr := net.UnixAddr{"/tmp/scpc", t}
	conn, err := net.DialUnix(t, &laddr, &net.UnixAddr{socketPath, t})
	if err != nil {
		panic(err)
	}

	sent := send(os.Args[1:])
	_, err = conn.Write(sent)
	if err != nil {
		panic(err)
	}

	rsp := response(conn, sent)
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
	flag.StringVar(&socketEnv, "env", socketEnv, "Read the socket from the given env variable")
	flag.StringVar(&socketPath, "path", defaulSocketPath(), "Read the socket from the given path.")
	flag.DurationVar(&responseTimeout, "timeout", responseTimeout, "Wait only specified seconds euclid's response, default is 2")
	flag.BoolVar(&verbose, "verbose", verbose, "Verbose logging messages, default is false.")
	flag.Parse()
}
