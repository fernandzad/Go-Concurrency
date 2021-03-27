package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func SimpleConcurrency() {
	// This adds a goroutine and keeps executing the main program
	go count("something")
	count("else")

	// This adds two goroutine and continue the main program until it exits
	// So this way there is no time to execute any goroutine
	go count("something")
	go count("else")
}

// When we need to wait for goroutines to finish in order to continue the program
func WaitGroupDemo() {
	// This creates a wait group that allows us to wait until our goroutines finish
	var wg sync.WaitGroup
	wg.Add(2) // We are adding two things to wait

	go func() {
		countToFive("something")
		wg.Done() // We are done with the first (we know is the first because countToFive has less operations than countToTen)
	}()

	go func() {
		countToTen("else")
		wg.Done() // We are done with the second
	}()

	wg.Wait() // We'll continue with the execution of the main until our group finishes

	fmt.Println("This is Concurrency!")
}

func ChannelsSelect() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			c1 <- "From Channel 1: Every .5 seconds"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		for {
			c2 <- "From Channel 2: Every 2 seconds"
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		select {
		case msgFromChannel1 := <-c1:
			fmt.Println(msgFromChannel1)
		case msgFromChannel2 := <-c2:
			fmt.Println(msgFromChannel2)
		}
	}
}

func ChannelDemo() {
	c := make(chan string)
	go countChannel("dig", c)

	// for {
	// 	msgFromChannel, isOpen := <-c
	// 	if !isOpen {
	// 		fmt.Println("Channel closed")
	// 		break
	// 	}
	// 	fmt.Println("Message from channel", msgFromChannel)
	// }
	for msgFromChannel := range c {
		fmt.Println("Message from channel", msgFromChannel)
	}
}

func countToTen(thing string) {
	for i := 1; i <= 10; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}

func count(thing string) {
	for i := 1; true; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}

func countChannel(thing string, c chan string) {
	for i := 1; i <= 5; i++ {
		c <- (thing + strconv.Itoa(i))
		time.Sleep(time.Millisecond * 500)
	}
	// The sender is the responsible to close the channel
	// The receiver never has to close the channel because
	// it can close the channel prematurely and this will throw a panic
	close(c)
}

func countToFive(thing string) {
	for i := 1; i <= 5; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}
