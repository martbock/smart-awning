package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"smart-awning/awning"
	"smart-awning/rest"
	"sync"
	"syscall"
)

func main() {
	var cancelFuncs []*context.CancelFunc
	setupCloseHandler(cancelFuncs)
	wg := sync.WaitGroup{}
	wg.Add(2)

	commandChan, cancel, err := awning.Prepare(&wg)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	cancelFuncs = append(cancelFuncs, cancel)

	cancel, err = rest.Start(commandChan, &wg)
	if err != nil {
		return
	}
	cancelFuncs = append(cancelFuncs, cancel)

	wg.Wait()

	log.Println("Shutting down now")
	os.Exit(0)
}

func setupCloseHandler(cancelFuncs []*context.CancelFunc) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Program was cancelled, stopping goroutines")
		for _, cancelFunc := range cancelFuncs {
			(*cancelFunc)()
		}
	}()
}
