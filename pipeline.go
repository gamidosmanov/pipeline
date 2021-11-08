package main

import "log"

type Pipeline struct {
	done   <-chan bool
	Stages []Stage
}

func NewPipeline(done <-chan bool) *Pipeline {
	stages := make([]Stage, 0)
	log.Println("Pipeline created")
	return &Pipeline{done: done, Stages: stages}
}

func (p *Pipeline) AddStage(s Stage) {
	p.Stages = append(p.Stages, s)
}

func (p *Pipeline) applyStage(s Stage, done <-chan bool, input <-chan int) <-chan int {
	return s(done, input)
}

func (p *Pipeline) Start(input <-chan int) <-chan int {
	var data <-chan int = input
	for i := range p.Stages {
		data = p.applyStage(p.Stages[i], p.done, data)
	}
	log.Println("Pipeline started")
	return data
}
