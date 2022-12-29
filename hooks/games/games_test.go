package games

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Log(r1.Intn(10))

	<-time.After(time.Second * 5)
	r2 := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Log(r2.Intn(10))
	t.Log(r1.Intn(10))
}

func TestMapEmpty(t *testing.T) {
	m := make(map[int]uint8)
	t.Log(m[1])
	m[1] += 30
	t.Log(m[1])
}

func TestFixedArr(t *testing.T) {
	var a [4]int
	t.Log(len(a))
	t.Log(a)
	a[1] = 1
	t.Log(a)
	t.Log(len(a))
}

func TestChar(t *testing.T) {
	points := uint8(0)
	card1, card2, card3 := "10♠", "3♠", "A♠"
	for _, v := range []string{card1, card2, card3} {
		// remove suit
		ru := []rune(v)
		v = string(ru[:len(ru)-1])
		switch v {
		case "A":
			points += 1
		case "J", "Q", "K", "10":
			points += 10
		default:
			points += uint8(v[0] - '0')
		}
	}
	t.Log(points)
}

func TestContextxDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*3))
	go func() {
		<-ctx.Done()
		defer cancel()
		t.Log("activated!")
	}()
	go func() {
		for i := 5; i > 0; i-- {
			<-time.After(time.Second)
			select {
			case <-ctx.Done():
				t.Log("done")
			default:
				t.Log(i)
				if i == 2 {
					cancel()
					t.Log("cancel")
				}
			}
		}
	}()

	<-time.After(time.Second * 10)
}
