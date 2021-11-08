package main

import (
	"fmt"
	"log"
	"time"
)

const (
	bufferSize   int           = 5
	flushTimeout time.Duration = 10 * time.Second
)

func main() {
	// Запускаем поток данных
	input, done := launchReader()
	log.Println("Reader started")

	// Создаем пайплайн и добавляем стадии обработки
	pipe := NewPipeline(done)
	pipe.AddStage(positiveStage)
	log.Println("Stage added: positive")
	pipe.AddStage(thirdsStage)
	log.Println("Stage added: thirds")
	pipe.AddStage(bufferingStage)
	log.Println("Stage added: buffering")

	// Запускаем пайплайн
	data := pipe.Start(input)

	// Потребитель
	go func(done <-chan bool, input <-chan int) {
		for {
			select {
			case <-done:
				return
			case i := <-input:
				fmt.Printf("Получены данные: %d\n", i)
			}
		}
	}(done, data)
	log.Println("Receiver started")

	// Ждем сигнала об окончании
	<-done

	// Ура, оно работает
}
