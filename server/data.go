package main

type classQueryInput struct {
	GenEds     *[]string
	MinCredits *int
	MaxCredits *int
}

type professor interface {
	Name() string
}
