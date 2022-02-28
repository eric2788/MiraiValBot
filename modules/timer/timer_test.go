package timer

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-ctx.Done()
		fmt.Println("done A")
	}()

	go func() {
		<-ctx.Done()
		fmt.Println("done B")
	}()

	<-time.After(time.Second * 2)
	cancel()
}

func TestPtrAssign(t *testing.T) {
	started := true
	assignBool(&started)
	assert.Equal(t, false, started)
}

func assignBool(b *bool) {
	*b = false
}

func ATestTicker(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	ticker := time.NewTicker(time.Minute * 2)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("timer stopped")
				return
			case tt := <-ticker.C:
				fmt.Println(tt)
			}
		}
	}()

	<-time.After(time.Second * 6)
	cancel()
	fmt.Println("stopped, waiting for done")
	wg.Wait()
}
