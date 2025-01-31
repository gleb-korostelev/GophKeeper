// Package closer provides a mechanism to manage and close external connections and resources gracefully.
//
// This package ensures that all external connections (e.g., database connections, network clients)
// are closed properly upon application termination. It also provides synchronization to manage
// concurrent closures.

package closer

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var globalCloser *closer

func init() {
	// Initialize a global closer instance that listens for termination signals (SIGTERM, SIGINT).
	globalCloser = New(syscall.SIGTERM, syscall.SIGINT)
}

// closer manages a list of resources to be closed and ensures they are closed
// gracefully when a termination signal is received or when explicitly triggered.
type closer struct {
	mu      sync.Mutex    // Protects access to the toClose list.
	done    chan struct{} // Signals when all resources have been closed.
	toClose []Closer      // List of resources to close.
	once    sync.Once     // Ensures closeAll is executed only once.
}

// Closer defines an interface that any resource needing cleanup must implement.
//
// Close must release the resource and return an error if it fails.
type Closer interface {
	Close() error
}

// New creates a new instance of `closer` and optionally starts listening for
// operating system signals.
func New(sig ...os.Signal) *closer {
	c := &closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.closeAll()
		}()
	}
	return c
}

// Add appends resources to the global closer's list for management.
func Add(conn ...Closer) {
	globalCloser.add(conn...)
}

// CloseAll closes all managed resources in the global closer.
//
// This function blocks until all resources have been closed.
func CloseAll() {
	globalCloser.closeAll()
}

// Wait blocks until all managed resources in the global closer have been closed.
func Wait() {
	globalCloser.wait()
}

// wait blocks until the `done` channel is closed, indicating all resources have been closed.
func (c *closer) wait() {
	<-c.done
}

// add appends resources to the closer's internal list.
func (c *closer) add(conn ...Closer) {
	c.mu.Lock()
	c.toClose = append(c.toClose, conn...)
	c.mu.Unlock()
}

// closeAll closes all managed resources in the closer. It ensures this operation
// is executed only once, even if called multiple times.
func (c *closer) closeAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.toClose
		c.toClose = nil
		c.mu.Unlock()

		var wg sync.WaitGroup
		for _, f := range funcs {
			wg.Add(1)
			go func(wg *sync.WaitGroup, f Closer) {
				_ = f.Close() // Ignore errors from Close.
				wg.Done()
			}(&wg, f)
		}
		wg.Wait()
	})
}
