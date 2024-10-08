package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ChrisGora/semaphore"
)

type buffer struct {
	b                 []int
	size, read, write int
}

func newBuffer(size int) buffer {
	return buffer{
		b:     make([]int, size),
		size:  size,
		read:  0,
		write: 0,
	}
}

func (buffer *buffer) get() int {
	x := buffer.b[buffer.read]
	fmt.Println("Get\t", x, "\t", buffer)
	buffer.read = (buffer.read + 1) % len(buffer.b)
	return x
}

func (buffer *buffer) put(x int) {
	buffer.b[buffer.write] = x
	fmt.Println("Put\t", x, "\t", buffer)
	buffer.write = (buffer.write + 1) % len(buffer.b)
}

func producer(buffer *buffer, spaceAvailable, workAvailable semaphore.Semaphore, mutex *sync.Mutex, start, delta int) {
	x := start
	for {
		spaceAvailable.Wait() // Wait until there is space available in the buffer
		mutex.Lock()          // Lock the buffer to ensure exclusive access
		buffer.put(x)
		workAvailable.Post() // Signal that work is avaible
		mutex.Unlock()       // Unlock the buffer to ensure other can access

		x = x + delta
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func consumer(buffer *buffer, spaceAvailable, workAvailable semaphore.Semaphore, mutex *sync.Mutex) {
	for {
		workAvailable.Wait()
		mutex.Lock()
		_ = buffer.get()
		spaceAvailable.Post()
		mutex.Unlock()

		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	}
}

func main() {
	// 2a
	// 4, 994, 5

	// 2c How might you improve the buffer methods
	// 1. Use RWMutex, so that there's mulitple read rountine, improve concurrency
	// 2. Condition Variables

	buffer := newBuffer(5)
	mutex := &sync.Mutex{} // Initialse a mutex
	// var mutex sync.Mutex

	spaceAvailable := semaphore.Init(5, 5) // Initialise a semaphore for spaces available
	workAvailable := semaphore.Init(5, 0)  // Initialise a semaphore for work available

	// Start two producer goroutines with different starting values and increments.
	go producer(&buffer, spaceAvailable, workAvailable, mutex, 1, 1)
	go producer(&buffer, spaceAvailable, workAvailable, mutex, 1000, -1)

	consumer(&buffer, spaceAvailable, workAvailable, mutex)
	// ${[contents] length start_index end_index}
}
