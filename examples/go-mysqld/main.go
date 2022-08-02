// Copyright (C) 2020 Satoshi Konno. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
 go-mysqld is an example of implementing a compatible MySQL server using go-mysql.
	NAME
	 go-mysqld

	SYNOPSIS
	 go-mysqld [OPTIONS]

	OPTIONS
	-v      : Enable verbose output.
	-p      : Enable profiling.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"flag"
	"log"
	"net/http"

	// nolint
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/go-mysql/examples/go-mysqld/server"
)

const (
	ProgramName = " go-mysqld"
)

func main() {
	// Command Line Options

	//verbose := flag.Bool("v", false, "Verbose log")
	isProfileEnabled := flag.Bool("p", false, "Enable profiling")
	flag.Parse()

	if *isProfileEnabled {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	// Start server

	server := server.NewServer()
	err := server.Start()
	if err != nil {
		log.Printf("%s couldn't be started (%s)", ProgramName, err.Error())
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM)

	exitCh := make(chan int)

	go func() {
		for {
			s := <-sigCh
			switch s {
			case syscall.SIGHUP:
				log.Printf("Caught SIGHUP, restarting...")
				err = server.Restart()
				if err != nil {
					log.Printf("%s couldn't be restarted (%s)", ProgramName, err.Error())
					os.Exit(1)
				}
			case syscall.SIGINT, syscall.SIGTERM:
				log.Printf("Caught %s, stopping...", s.String())
				err = server.Stop()
				if err != nil {
					log.Printf("%s couldn't be stopped (%s)", ProgramName, err.Error())
					os.Exit(1)
				}
				exitCh <- 0
			}
		}
	}()

	code := <-exitCh

	os.Exit(code)
}
