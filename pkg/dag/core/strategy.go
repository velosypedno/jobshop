package core

type Strategy interface {
	Plan(*Problem) Solution
	Type() string
	Name() string
	Description() string
}
