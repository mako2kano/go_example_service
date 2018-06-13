package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kardianos/service"
)

var logger service.Logger

type exarvice struct {
	exit chan struct{}
}

func (e *exarvice) run() error {

	logger.Info("Exarvice Start !!!")

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case tm := <-ticker.C:
			logger.Infof("Still running at %v", tm)
		case <-e.exit:
			ticker.Stop()
			logger.Info("Exarvice Stop ...")
			return nil
		}
	}
}

func (e *exarvice) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	e.exit = make(chan struct{})

	go e.run()
	return nil
}

func (e *exarvice) Stop(s service.Service) error {
	close(e.exit)
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "Exarvice",
		DisplayName: "Exarvice (Go Service Example)",
		Description: "This is an example Go service.",
	}

	// Create Exarvice service
	program := &exarvice{}
	s, err := service.New(program, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Setup the logger
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal()
	}

	if len(os.Args) > 1 {

		err = service.Control(s, os.Args[1])
		if err != nil {
			fmt.Printf("Failed (%s) : %s\n", os.Args[1], err)
			return
		}
		fmt.Printf("Succeeded (%s)\n", os.Args[1])
		return
	}

	// run in terminal
	s.Run()
}
