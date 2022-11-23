package main

/*
 * @notice this in command line tool for finding all primal number within input ranges
 * all numbers will be written in file with inputed name (default primes.txt)
 * writter and finder are gorutines, that connetct with each other by channels
 * all this proccess will be killed after deadline (by context)
 */

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var numbers = make(chan string)
var done = make(chan bool)
var exit = make(chan bool)

// custop flag type for multiple ranges
type ranges []int

func (r *ranges) String() string {
	return fmt.Sprint(*r)
}

func (r *ranges) Set(value string) error {
	for _, n := range strings.Split(value, ":") {
		num, err := strconv.Atoi(n)
		if err != nil {
			return err
		}
		*r = append(*r, num)
	}
	return nil
}

func main() {
	var ranges ranges
	filename := flag.String("file", "primes.txt", "Name of generated file with all finded numbers")
	timeout := flag.Int("timeout", 10, "Timeout for runing the program in seconds")
	flag.Var(&ranges, "range", "In which interval to look for prime numbers")
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(*timeout)*time.Second)
	defer cancel()

	if len(ranges) == 0 {
		fmt.Println("Nil range")
		return
	}

	go write(ctx, *filename, len(ranges)/2)
	for i := 0; i < len(ranges); i += 2 {
		if ranges[i] > ranges[i+1] {
			panic("Rigth edge less then left")
		}
		go findPrime(ranges[i], ranges[i+1])
	}
	select {
	case <-exit:
		fmt.Println("All done!")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}

/**
 * @notice function for finding primes
 */
func findPrime(from, to int) {
	var wg sync.WaitGroup
	if from%2 == 0 {
		from += 1
	}
	for i := from; i <= to; i += 2 {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			for j := 3; j < num; j += 2 {
				if j*(j-1) > num {
					break
				}
				if num%j == 0 {
					return
				}
			}
			numbers <- fmt.Sprintf("%v\n", num)
		}(i)
	}
	wg.Wait()
	// fmt.Printf("find %v:%v done\n", from, to)
	done <- true
}

/**
 * @notice write all finded numbers in file, one number in row
 * if file already exist then just add this numbers in the end of file
 */
func write(ctx context.Context, filename string, workers int) {
	ended := 0
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic("Can't open/create file")
	}
	defer f.Close()
	for {
		select {
		case n := <-numbers:
			f.WriteString(n)
		case <-done:
			ended++
			if ended == workers {
				exit <- true
				return
			}
		}
	}
}
