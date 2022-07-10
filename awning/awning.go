package awning

import (
	"context"
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"smart-awning/config"
	"sync"
	"time"
)

var (
	pinExtend  = rpio.Pin(config.GPIO.PinExtend)
	pinRetract = rpio.Pin(config.GPIO.PinRetract)
	pinStop    = rpio.Pin(config.GPIO.PinStop)
)

func Prepare(wg *sync.WaitGroup) (chan<- Command, *context.CancelFunc, error) {
	err := rpio.Open()
	if err != nil {
		return nil, nil, err
	}

	for _, pin := range []rpio.Pin{pinExtend, pinRetract, pinStop} {
		pin.Output()
		pin.High()
	}

	c := make(chan Command, 100)

	ctx, cancelFunc := context.WithCancel(context.Background())

	go loop(c, wg, ctx)

	return c, &cancelFunc, nil
}

func loop(c <-chan Command, wg *sync.WaitGroup, ctx context.Context) {
	defer func() {
		_ = rpio.Close()
		wg.Done()
	}()

	for {
		select {
		case command := <-c:
			switch command {
			case Extend:
				log.Println("Extending awning")
				trigger(pinExtend)
			case Retract:
				log.Println("Retracting awning")
				trigger(pinRetract)
			case Stop:
				log.Println("Stopping awning")
				trigger(pinStop)
			}
		case <-ctx.Done():
			return
		}
	}
}

func trigger(pin rpio.Pin) {
	pin.Low()
	time.Sleep(time.Millisecond * 300)
	pin.High()
}
