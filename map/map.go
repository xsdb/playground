package main

import (
	"fmt"
	"log"
	"sync"
)

type Value struct {
	value string
}

var (
	wg     sync.WaitGroup
	nrLoop int = 1000000
)

func prepare() map[string]*Value {
	m := make(map[string]*Value)
	for i := 0; i < 10; i++ {
		k := fmt.Sprintf("key-%v", i)
		v := fmt.Sprintf("value-%v", i)
		m[k] = &Value{v}
	}
	return m
}

func read(m map[string]*Value, key string) {
	defer wg.Done()
	for i := 0; i < nrLoop; i++ {
		log.Printf("read=%v\n", m[key])
	}
}

func write(m map[string]*Value, key string) {
	defer wg.Done()
	for i := 0; i < nrLoop; i++ {
		v := fmt.Sprintf("value-%v", i)
		m[key] = &Value{v}
	}
}

func update(m map[string]*Value, key string) {
	defer wg.Done()
	for i := 0; i < nrLoop; i++ {
		v := fmt.Sprintf("value-%v", i)
		m[key].value = v
	}
}

func readSameKeys(m map[string]*Value) {
	go read(m, "key-1")
	go read(m, "key-1")
}

func readDiffKeys(m map[string]*Value) {
	go read(m, "key-1")
	go read(m, "key-2")
}

func writeSameKeys(m map[string]*Value) {
	go write(m, "key-1")
	go write(m, "key-1")
}

func updateSameKeys(m map[string]*Value) {
	go update(m, "key-1")
	go update(m, "key-1")
}

func writeDiffKeys(m map[string]*Value) {
	go write(m, "key-1")
	go write(m, "key-2")
}

func updateDiffKeys(m map[string]*Value) {
	go update(m, "key-1")
	go update(m, "key-2")
}

func readWriteSameKeys(m map[string]*Value) {
	go read(m, "key-1")
	go write(m, "key-1")
}
func readUpdateSameKeys(m map[string]*Value) {
	go read(m, "key-1")
	go update(m, "key-1")
}

func readWriteDiffKeys(m map[string]*Value) {
	go read(m, "key-1")
	go write(m, "key-2")
}

func readUpdateDiffKeys(m map[string]*Value) {
	go read(m, "key-1")
	go update(m, "key-2")
}

func main() {
	m := prepare()
	wg.Add(2)

	readSameKeys(m) // => OK
	// readDiffKeys(m) // => OK
	// writeSameKeys(m) // => fatal error: concurrent map writes
	// writeDiffKeys(m) // => fatal error: concurrent map writes
	// readWriteSameKeys(m) // => fatal error: concurrent map read and map write
	// readWriteDiffKeys(m) // => fatal error: concurrent map read and map write

	// updateSameKeys(m) // => OK
	// updateDiffKeys(m) // => OK
	// readUpdateSameKeys(m) // => OK
	// readUpdateDiffKeys(m) // => OK

	wg.Wait()
	log.Println("Finished")
}
