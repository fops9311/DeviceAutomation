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
	Tank1.program = func(in Recipe, done <-chan interface{}) {
		fmt.Printf("Processing %v...\n", in.GetName())
		Delay := time.NewTimer(time.Second * 1)
		<-Delay.C
	}
	Tank2.program = func(in Recipe, done <-chan interface{}) {
		fmt.Printf("Processing %v...\n", in.GetName())
		Delay := time.NewTimer(time.Second * 1)
		<-Delay.C
	}
	Tank3.program = func(in Recipe, done <-chan interface{}) {
		fmt.Printf("Processing %v...\n", in.GetName())
		Delay := time.NewTimer(time.Second * 1)
		<-Delay.C
	}
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
	program func(in Recipe, done <-chan interface{})
}

func (h *Tank) init() {
	fmt.Println(h.message)
}

func (h *Tank) RecipePipeline(in <-chan Recipe, done <-chan interface{}) (<-chan Recipe, <-chan interface{}) {
	out := make(chan Recipe)
	outDone := make(chan interface{})
	go func() {
		for {
			select {
			case h.recipe = <-in:
				h.program(h.recipe, outDone)
				fmt.Printf("Recipe %v processed in %v\n", h.recipe.GetName(), h.name)
				out <- h.recipe
			case <-done:
				close(outDone)
				return
			}
		}
	}()
	return out, outDone
}
