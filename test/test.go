package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var mu sync.Mutex
var requestsDone int

func runClient(id, cycles int, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "127.0.0.1:64666")
	if err != nil {
		fmt.Println("client", id, "dial error:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for i := id * cycles; i < id*cycles+cycles; i++ {
		commands := []string{
			fmt.Sprintf("SET %d localhost\n", i),
			fmt.Sprintf("GET %d\n", i),
			fmt.Sprintf("DEL %d\n", i),
		}

		for _, cmd := range commands {
			_, err := conn.Write([]byte(cmd))
			if err != nil {
				fmt.Println("client", id, "write error:", err)
				return
			}

			n, err := conn.Read(buf)

			if err != nil {
				fmt.Println("client", id, "read error:", err)
				return
			}

			// Delete this line below if you want to print
			n++
			// Remove the comment below if you want to print
			//fmt.Printf("client %d: %s", id, string(buf[:n]))

			mu.Lock()
			requestsDone++
			mu.Unlock()
		}
	}
}

func main() {
	startedAt := time.Now()
	var wg sync.WaitGroup

	numClients := 10        // número de clientes concurrentes
	cyclesPerClient := 5000 // cada cliente hace 50 ciclos (SET, GET, DEL)

	for c := 0; c < numClients; c++ {
		wg.Add(1)
		go runClient(c, cyclesPerClient, &wg)
	}

	wg.Wait()

	elapsed := time.Since(startedAt).Seconds()
	fmt.Println(requestsDone, "requests done!")
	fmt.Printf("-----------------\n%.2f seconds\n%.2f req/s\n",
		elapsed, float64(requestsDone)/elapsed)
}
