package kivaio

import (
	"regexp"
	"strconv"

	"github.com/st3v/tracerr"
)

type Parser interface {
	Parse(message string) (*message, error)
}

type parser struct{}

func newParser() *parser {
	return &parser{}
}

func (p *parser) Parse(data string) (*message, error) {
	r := regexp.MustCompile("([^:]+):([0-9]+)?(\\+)?:([^:]+)?:?([\\s\\S]*)?")
	parts := r.FindStringSubmatch(data)

	category, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &message{
		id:       parts[2],
		category: category,
		data:     parts[5],
		endpoint: parts[4],
		ack:      parts[3],
	}, nil
}
