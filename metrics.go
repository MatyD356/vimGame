package main

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
)

func countFuncExecution(f func(), wg *sync.WaitGroup) {
	defer wg.Done()
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	start := time.Now()
	f()
	elapsed := time.Since(start).Seconds()
	fmt.Printf("Function '%s' executed in: %.3f seconds\n", funcName, elapsed)
}
