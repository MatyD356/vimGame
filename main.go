package main

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type Player struct {
	X, Y int
}

func countFuncExecution(f func(), wg *sync.WaitGroup) {
	defer wg.Done()
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	start := time.Now()
	f()
	elapsed := time.Since(start).Seconds()
	fmt.Printf("Function '%s' executed in: %.3f seconds\n", funcName, elapsed)
}

func twoSeconds() {
	time.Sleep(2123 * time.Millisecond)
}

func fiveSeconds() {
	time.Sleep(5 * time.Second)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go countFuncExecution(twoSeconds, &wg)
	go countFuncExecution(fiveSeconds, &wg)
	wg.Wait()
}
