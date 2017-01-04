package main

import (
	"fmt"
	"log"

	"github.com/sethgrid/chring"
)

func main() {
	rm := chring.NewRingManager()
	for _, n := range []string{"123.45.83.190", "123.45.83.191", "123.45.83.192", "123.45.78.191", "123.45.78.189", "123.12.09.249"} {
		rm.AddNode(n)
	}
	for i := 1; i <= 100; i++ {
		rm.AddKey(fmt.Sprintf("user_%d", i))
	}
	log.Println("open http://locahost:5000")
	chring.ServeRingManager(rm, ":5000")
}
