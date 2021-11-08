package main

import (
	"fmt"
	"time"
)

const (
	bufferSize   int           = 5
	flushTimeout time.Duration = 10 * time.Second
)

func main() {
	// Запускаем поток данных
	input, done := launchReader()

	// Создаем пайплайн и добавляем стадии обработки
	pipe := NewPipeline(done)
	pipe.AddStage(positiveStage)
	pipe.AddStage(thirdsStage)
	pipe.AddStage(bufferingStage)

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

	// Ждем сигнала об окончании
	<-done

	// Ура, оно работает
}
