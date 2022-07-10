package awning

type Command int

const (
	Extend Command = iota
	Retract
	Stop
)
