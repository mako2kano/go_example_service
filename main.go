package main

import (
	"errors"
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

func (e *exarvice) Run() error {

	logger.Info("Exarvice Start !!!")

	ticker := time.NewTicker(5 * time.Second)
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

	go e.Run()
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
		switch os.Args[1] {
		case "install":
			err = s.Install()
		case "uninstall":
			err = s.Uninstall()
		case "start":
			err = s.Start()
		case "stop":
			err = s.Stop()
		case "run":
			err = s.Run()
		default:
			err = errors.New("invalid argument")
		}
		if err != nil {
			fmt.Printf("Failed (%s) : %s\n", os.Args[1], err)
			return
		}
		fmt.Printf("Succeeded (%s)\n", os.Args[1])
		return
	}

	fmt.Println("Usage: go_example_service [OPTION]")
	fmt.Println("  install  : install serve")
	fmt.Println("  uninstall: uninstall service")
	fmt.Println("  start    : start service program")
	fmt.Println("  stop     : stop service program")
	fmt.Println("  run      : run in terminal")
}
