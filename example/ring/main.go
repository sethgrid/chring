package main

import (
	"log"

	"github.com/sethgrid/chring"
)

func main() {
	ring := chring.NewRing()
	for _, n := range []string{"123.45.83.190", "123.45.83.191", "123.45.83.192", "123.45.78.191", "123.45.78.189", "123.12.09.249"} {
		ring.Add(n)
	}
	log.Println("open http://locahost:5000")
	chring.ServeRing(ring, ":5000")
}
