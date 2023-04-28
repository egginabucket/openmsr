package libtracks

import (
	"fmt"
	"strings"
)

type (
	Info struct {
		name     string
		value    any
		branches []Informer
	}
	Informer interface {
		Info() *Info
	}
)

func (i *Info) Info() *Info {
	// could be a second type lol
	return i
}

func (i *Info) write(b *strings.Builder, depth int) {
	if depth > 0 {
		b.WriteByte('\n')
		for d := 0; d < depth; d++ {
			b.WriteByte('\t')
		}
	}
	fmt.Fprintf(b, "%s: %v", i.name, i.value)
	for _, ni := range i.branches {
		ni.Info().write(b, depth+1)
	}
}

func (i *Info) String() string {
	var b strings.Builder
	i.write(&b, 0)
	return b.String()
}

func NewInfo(name string, value any, branches ...Informer) *Info {
	return &Info{
		name:     name,
		value:    value,
		branches: branches,
	}
}
