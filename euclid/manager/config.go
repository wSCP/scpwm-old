package manager

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

const (
	COMMENT      = byte('#')
	LINECONTINUE = byte('\\')
)

func (m *Manager) LoadConfig(f string) error {
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	if err != nil {
		return err
	}
	conf, err := parseConfig(reader)
	if err != nil {
		return err
	}
	go execConfig(conf)
	return nil
}

func parseConfig(r *bufio.Reader) ([]string, error) {
	lineno := 0
	var err error
	var ret []string
	for err == nil {
		l, _, err := r.ReadLine()
		if err != nil {
			break
		}
		lineno++
		if len(l) == 0 || l[0] == COMMENT {
			continue
		}
		for l[len(l)-1] == LINECONTINUE {
			nl, _, err := r.ReadLine()
			if err != nil {
				break
			}
			l = append(l, bytes.TrimFunc(nl, unicode.IsSpace)...)
		}
		ret = append(ret, string(l))
	}
	return ret, err
}

func execConfig(cmd []string) {
	for _, c := range cmd {
		spl := strings.Split(c, " ")
		rc := exec.Command(spl[0], spl[1:]...)
		go execCmd(rc)
	}
}

func execCmd(c *exec.Cmd) {
	c.Run()
}
