package query

import (
	"github.com/thrisp/scpwm/euclid/client"
	"github.com/thrisp/scpwm/euclid/monitor"
)

type coordinate struct {
	m monitor.Monitor
	c client.Client
}

var NoCoordinate = coordinate{nil, nil}

func Coordinate(m Monitor, c Client) coordinate {
	return coordinate{
		m: m,
		c: c,
	}
}
