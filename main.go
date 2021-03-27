package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"sync"
)

type Pokemon struct {
	ID   int    `json:ID`
	Name string `json:Name`
	URL  string `json:Url`
}

func ReadConcurrently(f *os.File) *[]Pokemon {

	reader := csv.NewReader(f)
	reader.Comma = ','
	reader.Comment = '#'
	reader.FieldsPerRecord = -1

	defer f.Close()
	tempPokemon := []Pokemon{}

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("There was something wrong %v\n", err.Error())
		}

		if line[0] != "" {
			id, err := strconv.Atoi(line[0])

			if err != nil {
				fmt.Printf("There was something wrong %v\n", err.Error())
			}

			tempPokemon = append(tempPokemon, Pokemon{
				ID:   id,
				Name: line[1],
				URL:  line[2],
			})
		}

	}
	return &tempPokemon
}

func main() {
	f, err := os.Open("./input.csv")
	if err != nil {
		fmt.Printf("There was something wrong %v\n", err.Error())
	}
	pokemons := ReadConcurrently(f)
	WorkerPools(1, 1, "even", *pokemons)
	WorkerPools(3, 2, "even", *pokemons)
	WorkerPools(5, 1, "even", *pokemons)
	WorkerPools(5, 2, "even", *pokemons)
	WorkerPools(6, 1, "even", *pokemons)
	WorkerPools(10, 3, "even", *pokemons)
}

func calculatePoolSize(items int, itemsPerWorker int, totalPokemons int) int {
	var poolSize int
	if items%itemsPerWorker != 0 {
		poolSize = int(math.Ceil(float64(items) / float64(itemsPerWorker)))
	} else {
		poolSize = int(items / itemsPerWorker)
	}

	// If we overpass the number of workers above the half of number
	// of items it's gonna get into an infinit looop
	if poolSize > (totalPokemons / 2) {
		poolSize = totalPokemons / 2
	}
	return poolSize
}

func calculateMaxPokemons(totalPokemons int) int {
	var maxPokemons int

	if totalPokemons%2 == 0 {
		maxPokemons = totalPokemons / 2
	} else {
		maxPokemons = totalPokemons/2 + 1
	}
	return maxPokemons
}

func WorkerPools(items int, itemsPerWorker int, typeNumber string, pokemons []Pokemon) {
	totalPokemons := len(pokemons)
	poolSize := calculatePoolSize(items, itemsPerWorker, totalPokemons)
	maxPokemons := calculateMaxPokemons(totalPokemons)

	values := make(chan int)
	jobs := make(chan int, poolSize)
	shutdown := make(chan struct{})

	startIndex := 0
	var limit int
	limit = int(math.Ceil(float64(totalPokemons) / float64(poolSize)))
	lastLimit := (totalPokemons % limit)

	var wg sync.WaitGroup
	wg.Add(poolSize)

	for i := 0; i < poolSize; i++ {
		go func(jobs <-chan int) {
			for {
				var id int
				var limitRecalculated int
				start := <-jobs

				if limit+start >= totalPokemons && lastLimit != 0 {
					limitRecalculated = start + lastLimit
				} else {
					limitRecalculated = start + limit
				}

				for j := start; j < limitRecalculated; j++ {
					id = pokemons[j].ID

					select {
					case values <- id:
					case <-shutdown:
						wg.Done()
						return
					}
				}

			}
		}(jobs)
	}

	for i := 0; i < poolSize; i++ {
		jobs <- startIndex
		startIndex += limit
	}
	close(jobs)

	var nums []Pokemon = nil
	bucket := make(map[int]int, totalPokemons+1)
	for elem := range values {
		if elem%2 != 0 && bucket[elem] == 0 {
			nums = append(nums, pokemons[elem-1])
			bucket[elem] = elem
		}
		if len(nums) >= items || len(nums) >= maxPokemons {
			break
		}
	}

	fmt.Println("Receiver slimiting shutdown signal")
	close(shutdown)

	wg.Wait()

	fmt.Printf("Result count: %d\n", len(nums))
	fmt.Println(nums)
}
