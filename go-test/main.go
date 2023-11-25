package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/montanaflynn/stats"
	"github.com/redis/go-redis/v9"
)

const iterationsNumber = 1000000

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func BenchmarkRedisRDBSet() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	// disable AOF
	_, err = rdb.ConfigSet(ctx, "appendonly", "no").Result()
	if err != nil {
		panic(err)
	}

	// assure RDB is set to save snapshots every one second if at least one record exist
	_, err = rdb.ConfigSet(ctx, "save", "1 10").Result()
	if err != nil {
		panic(err)
	}

	timings := make([]float64, 0)
	for i := 0; i < iterationsNumber; i++ {
		started := time.Now()
		_, err := rdb.Set(ctx, randSeq(5), "value", 0).Result()
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(started)
		timings = append(timings, float64(elapsed.Nanoseconds()))
	}

	mean, _ := stats.Mean(timings)
	p90th, _ := stats.Percentile(timings, 90)
	p99th, _ := stats.Percentile(timings, 99)

	fmt.Println("============================================")
	fmt.Printf("Mean ----->  %f nanos\n", mean)
	fmt.Printf("(90 Percentile) ----->  %f nanos\n", p90th)
	fmt.Printf("(99 Percentile) ----->  %f nanos\n", p99th)
}

func BenchmarkRedisAOFBSet() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	// enable AOF
	_, err = rdb.ConfigSet(ctx, "appendonly", "yes").Result()
	if err != nil {
		panic(err)
	}

	// assure RDB is off
	_, err = rdb.ConfigSet(ctx, "save", "").Result()
	if err != nil {
		panic(err)
	}

	timings := make([]float64, 0)
	for i := 0; i < iterationsNumber; i++ {
		started := time.Now()
		_, err := rdb.Set(ctx, randSeq(5), "value", 0).Result()
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(started)
		timings = append(timings, float64(elapsed.Nanoseconds()))
	}

	mean, _ := stats.Mean(timings)
	p90th, _ := stats.Percentile(timings, 90)
	p99th, _ := stats.Percentile(timings, 99)

	fmt.Println("============================================")
	fmt.Printf("Mean ----->  %f nanos\n", mean)
	fmt.Printf("(90 Percentile) ----->  %f nanos\n", p90th)
	fmt.Printf("(99 Percentile) ----->  %f nanos\n", p99th)
}

func BenchmarkBeanstalkdPut() {
	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		panic(err)
	}

	timings := make([]float64, 0)
	for i := 0; i < iterationsNumber; i++ {
		started := time.Now()
		_, err := c.Put([]byte(randSeq(5)), 1, 0, 120*time.Second)
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(started)
		timings = append(timings, float64(elapsed.Nanoseconds()))
	}

	mean, _ := stats.Mean(timings)
	p90th, _ := stats.Percentile(timings, 90)
	p99th, _ := stats.Percentile(timings, 99)

	fmt.Println("============================================")
	fmt.Printf("Mean ----->  %f nanos\n", mean)
	fmt.Printf("(90 Percentile) ----->  %f nanos\n", p90th)
	fmt.Printf("(99 Percentile) ----->  %f nanos\n", p99th)
}

func BenchmarkBeanstalkdReserve() {
	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		panic(err)
	}

	timings := make([]float64, 0)
	for i := 0; i < iterationsNumber; i++ {
		started := time.Now()
		_, _, err := c.Reserve(0)
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(started)
		timings = append(timings, float64(elapsed.Nanoseconds()))
	}

	mean, _ := stats.Mean(timings)
	p90th, _ := stats.Percentile(timings, 90)
	p99th, _ := stats.Percentile(timings, 99)

	fmt.Println("============================================")
	fmt.Printf("Mean ----->  %f nanos\n", mean)
	fmt.Printf("(90 Percentile) ----->  %f nanos\n", p90th)
	fmt.Printf("(99 Percentile) ----->  %f nanos\n", p99th)
}

func main() {
	fmt.Printf("Benchmarking Beanstalkd PUT -- %d(iterations)\n", iterationsNumber)
	BenchmarkBeanstalkdPut()

	fmt.Printf("Benchmarking Beanstalkd RESERVE -- %d(iterations)\n", iterationsNumber)
	BenchmarkBeanstalkdReserve()

	fmt.Printf("Benchmarking Redis SET with RDB -- %d(iterations)\n", iterationsNumber)
	BenchmarkRedisRDBSet()

	fmt.Printf("Benchmarking Redis SET with AOF -- %d(iterations)\n", iterationsNumber)
	BenchmarkRedisAOFBSet()
}
