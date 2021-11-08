package main

import (
	"container/ring"
	"log"
	"sync"
)

type intBuffer struct {
	data *ring.Ring
	mu   *sync.Mutex
}

func newIntBuffer(size int) *intBuffer {
	data := ring.New(size)
	log.Println("Buffer created")
	return &intBuffer{data: data, mu: &sync.Mutex{}}
}

func (b *intBuffer) clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := 0; i < b.data.Len(); i++ {
		b.data.Value = nil
		b.data = b.data.Next()
	}
}

func (b *intBuffer) insert(i int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data.Value = i
	b.data = b.data.Next()
}

func (b *intBuffer) flush() (output []int) {
	b.mu.Lock()
	b.data.Do(func(i interface{}) {
		if i != nil {
			output = append(output, i.(int))
		}
	})
	b.mu.Unlock()
	b.clear()
	return
}
