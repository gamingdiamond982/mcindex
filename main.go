package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/dreamscached/minequery/v2"
	"github.com/zan8in/masscan"
)

var wg sync.WaitGroup

func parseResults(rstream bufio.Scanner) {
	defer wg.Done()
	for rstream.Scan() {
		srs := masscan.ParseResult(rstream.Bytes())
		fmt.Println("Proccessing: ", srs.IP, srs.Port)
		port, _ := strconv.Atoi(srs.Port)
		res, err := minequery.Ping17(srs.IP, port)
		if err != nil {
			log.Printf("Issue with %s:%s - %v\n", srs.IP, srs.Port, err)
			continue
		}
		fmt.Println(res)
	}
}

func main() {
	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets("192.168.1.0/24"),
		masscan.SetParamPorts("25565"),
		masscan.EnableDebug(),
		masscan.SetParamWait(0),
		masscan.SetParamRate(10000),
	)

	if err != nil {
		log.Fatalf("Unable to create masscan scanner: %v", err)
	}

	err = scanner.RunAsync()

	if err != nil {
		log.Fatalf("masscan encountered an error: %v", err)
	}

	stdout := scanner.GetStdout()
	stderr := scanner.GetStderr()

	wg.Add(2)
	go parseResults(stdout)
	go func() {
		defer wg.Done()
		for stderr.Scan() {
			fmt.Println(stderr.Text())
		}
	}()
	wg.Wait()
}
