package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func launchReader() (<-chan int, <-chan bool) {
	input := make(chan int)
	done := make(chan bool)
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(os.Stdin)
		var text string
		fmt.Println("Начните вводить целые числа или \"exit\" чтобы выйти")
		for {
			scanner.Scan()
			text = scanner.Text()
			if text == "exit" {
				log.Println("Input: exit")
				fmt.Println("До встречи!")
				return
			}
			i, err := strconv.Atoi(text)
			if err != nil {
				log.Println("Input: incorrect")
				fmt.Println("Нужно ввести целое число")
				continue
			}
			input <- i
			log.Println("Input: received")
		}
	}()
	return input, done
}
