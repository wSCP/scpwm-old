package version

import "fmt"

type Version interface {
	Package() string
	Tag() string
	Hash() string
	Date() string
	Fmt() string
}

func New(p, t, h, d string) Version {
	return &version{p, t, h, d}
}

type version struct {
	p, t, h, d string
}

func (v *version) Package() string {
	return v.p
}

func (v *version) Tag() string {
	return v.t
}

func (v *version) Hash() string {
	return v.h
}

func (v *version) Date() string {
	return v.d
}

func (v *version) Fmt() string {
	var msg string
	if v.h != "" && v.d != "" {
		msg = fmt.Sprintf("%s version %s(%s %s)\n", v.p, v.t, v.h, v.d)
	} else {
		msg = fmt.Sprintf("%s %s\n", v.p, v.t)
	}
	return msg
}
