package main

import "time"

func positiveStage(done <-chan bool, input <-chan int) <-chan int {
	positiveStream := make(chan int)
	go func() {
		defer close(positiveStream)
		for {
			select {
			case <-done:
				return
			case i := <-input:
				if i >= 0 {
					select {
					case positiveStream <- i:
					case <-done:
						return
					}
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
				if i != 0 && i%3 == 0 {
					select {
					case thirdsStream <- i:
					case <-done:
						return
					}
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
				for _, r := range results {
					bufferedStream <- r
				}
			}
		}
	}(ticker)
	return bufferedStream
}
