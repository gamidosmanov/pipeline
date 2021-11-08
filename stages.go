package main

import (
	"log"
	"time"
)

type Stage func(<-chan bool, <-chan int) <-chan int

func positiveStage(done <-chan bool, input <-chan int) <-chan int {
	positiveStream := make(chan int)
	go func() {
		defer close(positiveStream)
		for {
			select {
			case <-done:
				return
			case i := <-input:
				log.Printf("Positive stage: got %d\n", i)
				if i >= 0 {
					select {
					case positiveStream <- i:
						log.Printf("Positive stage: passed %d\n", i)
					case <-done:
						return
					}
				} else {
					log.Printf("Positive stage: rejected %d\n", i)
				}
			}
		}
	}()
	return positiveStream
}

func thirdsStage(done <-chan bool, input <-chan int) <-chan int {
	thirdsStream := make(chan int)
	go func() {
		defer close(thirdsStream)
		for {
			select {
			case <-done:
				return
			case i := <-input:
				log.Printf("Thirds stage: got %d\n", i)
				if i != 0 && i%3 == 0 {
					select {
					case thirdsStream <- i:
						log.Printf("Thirds stage: passed %d\n", i)
					case <-done:
						return
					}
				} else {
					log.Printf("Thirds stage: rejected %d\n", i)
				}
			}
		}
	}()
	return thirdsStream
}

func bufferingStage(done <-chan bool, input <-chan int) <-chan int {
	bufferedStream := make(chan int)
	buff := newIntBuffer(bufferSize)
	ticker := time.NewTicker(flushTimeout)
	// Буферизация
	go func(*intBuffer) {
		defer close(bufferedStream)
		for {
			select {
			case <-done:
				return
			case i := <-input:
				log.Printf("Buffering stage: got %d", i)
				buff.insert(i)
			}
		}
	}(buff)
	// Вывод буфера
	go func(*time.Ticker) {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				results := buff.flush()
				log.Println("Buffer flushed. Results obtained")
				for _, r := range results {
					bufferedStream <- r
				}
			}
		}
	}(ticker)
	return bufferedStream
}
