package broadcaster

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func ABenchmarkChannel(b *testing.B) {

	channel := make(chan int, 1024)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for i := 0; i < b.N; i++ {
			fmt.Printf("sent: %v\n", i)
			channel <- i
		}
		//close(channel)
		channel <- -1
	}()

	go func() {
		defer cancel()
		defer fmt.Println("Done.")
		for {
			s, ok := <-channel
			if !ok || s == -1 {
				return
			}
			fmt.Printf("received: %v\n", s)
			<-time.After(time.Second / 10)
		}
	}()

	<-ctx.Done()
}
