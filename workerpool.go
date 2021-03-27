package main

import "fmt"

func RunWorker() {
	// buffer of the channel is 100
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	go worker(jobs, results)
	go worker(jobs, results)
	go worker(jobs, results)
	go worker(jobs, results)

	for i := 0; i < 100; i++ {
		jobs <- i
	}
	close(jobs)

	for i := 0; i < 100; i++ {
		fmt.Println(<-results)
	}
}

func fibo(n int) int {
	if n <= 1 {
		return n
	}
	return fibo(n-1) + fibo(n-2)
}

func worker(jobs <-chan int, results chan<- int) {
	for n := range jobs {
		results <- fibo(n)
	}
}
