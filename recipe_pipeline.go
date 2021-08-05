package main

type Device interface {
	init()
	RecipePipeline(in <-chan Recipe, done <-chan interface{}) (<-chan Recipe, <-chan interface{})
}

type Recipe interface {
	GetName() string
}
