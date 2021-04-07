package value

import "strconv"

//go:generate stringer -type=Type -trimprefix=Type
type Type uint8

const (
	TypeNumber Type = iota
	TypeBoolean
	TypeChar
	TypeString
)

type Value interface {
	Type() Type
	String() string
}

type (
	Number struct {
		Val int
	}
	Boolean struct {
		Val bool
	}
	Char struct {
		Val rune
	}
	String struct {
		Val string
	}
)

func (*Number) Type() Type  { return TypeNumber }
func (*Boolean) Type() Type { return TypeBoolean }
func (*Char) Type() Type    { return TypeChar }
func (*String) Type() Type  { return TypeString }

func (n *Number) String() string {
	return strconv.Itoa(n.Val)
}

func (b *Boolean) String() string {
	if b.Val {
		return "true"
	} else {
		return "false"
	}
}

func (c *Char) String() string {
	return string(c.Val)
}

func (s *String) String() string {
	return s.Val
}
