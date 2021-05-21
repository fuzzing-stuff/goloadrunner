package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/fuzzing-stuff/chainer"
)

var (
	version = ""
)

func main() {

	nThreads := flag.Int("n", 1, "Number of simultaneous threads")
	inputFile := flag.String("f", "", "Input file")
	connectionString := flag.String("c", "", "Connection string")
	verInfo := flag.Bool("v", false, "Print version information and exit")

	flag.Parse()

	if *verInfo {
		fmt.Println("goloadrunner version ", version)
		return
	}

	reader, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	chains, err := chainer.LoadChains(reader)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(*nThreads)
	for i := 0; i < *nThreads; i++ {
		go func() {
			defer wg.Done()
			fmt.Print("+")
			conn, err := net.Dial("tcp", *connectionString)
			if err != nil {
				log.Println("ERROR connection to", *connectionString)
				return
			}
			for j := range chains {
				if msg, err := chains[j].Marshal(); err == nil {
					_, err := conn.Write(msg)
					if err != nil {
						log.Println("ERROR sending", *connectionString)
						return
					}

				} else {
					log.Println("ERROR marshaling item ", j)
					return
				}
				fmt.Print(".")
			}
			conn.Close()
		}()
	}
	//	ctx, cancel := context.WithCancel(context.Background())
	wg.Wait()
}
