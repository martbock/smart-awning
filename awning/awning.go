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
		pin.Low()
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
				extend()
			case Retract:
				retract()
			case Stop:
				stop()
			}
		case <-ctx.Done():
			return
		}
	}
}

func retract() {
	log.Println("Retracting awning")
	pinRetract.High()
	time.Sleep(time.Second * 1)
	pinRetract.Low()
}

func extend() {
	log.Println("Extending awning")
	pinExtend.High()
	time.Sleep(time.Second * 1)
	pinExtend.Low()
}

func stop() {
	log.Println("Stopping awning")
	pinStop.High()
	time.Sleep(time.Second * 1)
	pinStop.Low()
}
