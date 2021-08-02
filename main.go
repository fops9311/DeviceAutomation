package main

import (
	"fmt"
	"time"
)

func main() {
	Tank1 := Tank{
		name:    "T1",
		message: "this is Tank1",
	}
	Tank2 := Tank{
		name:    "T2",
		message: "this is Tank2",
	}
	Tank3 := Tank{
		name:    "T3",
		message: "this is Tank3",
	}
	Tank1.init()
	Tank2.init()
	Tank3.init()
	ProductionPlan := []BeerRecipe{
		{
			name:   "Cheshskoe",
			volume: 10,
		},
		{
			name:   "Hadizhinskoe",
			volume: 11,
		},
		{
			name:   "Ginnes",
			volume: 12,
		},
		{
			name:   "Ginnes 2",
			volume: 12,
		},
		{
			name:   "Cheshskoe",
			volume: 10,
		},
		{
			name:   "Hadizhinskoe",
			volume: 11,
		},
		{
			name:   "Ginnes 3",
			volume: 12,
		},
		{
			name:   "Ginnes 4",
			volume: 12,
		},
	}
	done := make(chan interface{})
	RecipeSender := func(ProductionPlan *[]BeerRecipe, done chan interface{}) (<-chan Recipe, <-chan interface{}) {
		out := make(chan Recipe)
		outDone := make(chan interface{})
		go func() {
			for _, Recipe := range *ProductionPlan {
				select {
				case out <- Recipe:
					fmt.Printf("Recipe %v added to the beggining of the pipeline\n", Recipe.GetName())
				case <-done:
					close(outDone)
					return
				}
			}
			close(outDone)
		}()
		return out, outDone
	}
	RecipeReciever := func(in <-chan Recipe, done <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})
		go func() {
			for {
				select {
				case Recipe := <-in:
					fmt.Printf("Recipe %v reached end of pipeline\n", Recipe.GetName())
				case <-done:
					close(out)
					return
				}
			}
		}()
		return out
	}
	Program1 := &testProgram{
		name: "Prog1",
	}
	Program2 := &testProgram{
		name: "Prog2",
	}
	Tank1.program = Program1
	Tank2.program = Program2
	Tank3.program = Tank1.program

	startTime := time.Now()
	<-RecipeReciever(
		Tank3.RecipePipeline(
			Tank2.RecipePipeline(
				Tank1.RecipePipeline(
					RecipeSender(&ProductionPlan, done),
				),
			),
		),
	)
	fmt.Println(time.Now().Sub(startTime))
	close(done)
}

type Recipe interface {
	GetName() string
}

type BeerRecipe struct {
	name   string
	volume float32
}

func (h BeerRecipe) GetName() string {
	switch h.name {
	default:
		return h.name
	}
}

type Device interface {
	init()
	RecipePipeline(in <-chan Recipe, done <-chan interface{}) (<-chan Recipe, <-chan interface{})
}
type Tank struct {
	name    string
	message string
	state   int
	recipe  Recipe
	program Program
}

func (h *Tank) init() {
	fmt.Println(h.message)
}

func (h *Tank) RecipePipeline(in <-chan Recipe, done <-chan interface{}) (<-chan Recipe, <-chan interface{}) {
	out := make(chan Recipe)
	outDone := make(chan interface{})

	progIn := make(chan Recipe)
	progInDone := make(chan interface{})

	go func() {
		for {
			select {
			case h.recipe = <-in:
				fmt.Printf("Recipe %v processed in %v\n", h.recipe.GetName(), h.name)

				progOut, progDone := h.program.Run(progIn, progInDone)

				progIn <- h.recipe

				out <- <-progOut
				<-progDone
			case <-done:
				close(progInDone)
				close(outDone)
				return
			}
		}
	}()
	return out, outDone
}

type Program interface {
	Run(in <-chan Recipe, done <-chan interface{}) (<-chan Recipe, <-chan interface{})
}

type testProgram struct {
	name string
}

func (h *testProgram) Run(in <-chan Recipe, done <-chan interface{}) (<-chan Recipe, <-chan interface{}) {
	out := make(chan Recipe)
	outDone := make(chan interface{})
	go func() {
		select {
		case recipe := <-in:
			fmt.Printf("%v Processing %v...\n", h.name, recipe.GetName())
			Delay := time.NewTimer(time.Second * 1)
			<-Delay.C
			out <- recipe
			fmt.Printf("Processing with %v of %v... ended\n", h.name, recipe.GetName())
			close(outDone)
			return
		case <-done:
			close(outDone)
			return
		}
	}()
	return out, outDone
}
