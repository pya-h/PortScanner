package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
)

func WorkerPool(address *string, ports chan uint16, results chan uint16, progress_trigger chan bool) {
	for port := range ports {
		if connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *address, port)); err == nil {
			results <- port
			connection.Close()
		} else {
			results <- 0
		}
		progress_trigger <- true
	}

}

func UpdateProgress(number_of_ports uint16, trigger chan bool) {
	checked_ports := uint16(0)
	var triggered bool
	progress_string := ""
	var progress float32
	var next_bar_point float32 = 2.5
	for {
		triggered = <-trigger
		if triggered {
			checked_ports++
		}
		progress = 100 * float32(checked_ports) / float32(number_of_ports)
		if progress > 100 {
			progress = 100
			progress_string = "\b"
		} else if progress >= next_bar_point {
			progress_string = "="
			next_bar_point += 2.5
		} else {
			progress_string = ""
		}
		fmt.Printf("\b\b\b\b\b\b\b\b%s>%5.2f %%", progress_string, progress)
	}
}

func main() {
	address := flag.String("h", "localhost", "Host address")
	first_port := flag.Int("first", 1, "First port")
	last_port := flag.Int("last", 1024, "Last port")
	workers_count := flag.Int("workers", 300, "Number of workers")

	flag.Parse()
	fmt.Println("Search started on:", *address, ", Ports:", *first_port, "->", *last_port)
	// The higher the count of workers(300 by default), the faster your program should execute. But if you add too many
	//workers, your results could become unreliable
	ports := make(chan uint16, *workers_count)
	results := make(chan uint16)
	// start goroutines waiting for channel
	capacity := int8(cap(ports))
	progress_trigger := make(chan bool)
	go UpdateProgress(uint16(*last_port-*first_port), progress_trigger)

	fmt.Print("\t  0.00 %")
	for i := int8(0); i < capacity; i++ {
		go WorkerPool(address, ports, results, progress_trigger) // create 100 go routines, which their execution is blocked on for loop,
		// until something is sent to the channel
	}

	go func() {
		// this must be go routine so that results channel can be handled in the next lines
		for i := *first_port; i <= *last_port; i++ {
			// sending data to then channel
			ports <- uint16(i)
		}
	}()

	var open_ports []uint16
	for i := *first_port; i <= *last_port; i++ {
		port := <-results
		if port > 0 {
			open_ports = append(open_ports, port)
		}

	}
	close(results)
	close(ports)

	// sort.Ints(open_ports as  int)
	sort.Slice(open_ports, func(i, j int) bool { return open_ports[i] < open_ports[j] })
	open_ports_count := len(open_ports)
	fmt.Print("\n\nOpen Ports: [")
	for index, port := range open_ports {
		fmt.Printf("%d", port)
		if next_index := index + 1; next_index < open_ports_count {
			if next_index%5 == 0 {
				fmt.Println()
			} else {
				fmt.Print("\t")
			}
		}
	}
	fmt.Println("]")
}
