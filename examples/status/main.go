package main

import (
	"context"
	"fmt"
	"log"

	mosquitto "example.com/mosquitto-manager"
)

func main() {
	b := mosquitto.New(mosquitto.DefaultConfig())

	status, err := b.Status(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("service=%s running=%v pid=%d\n", status.Service, status.Running, status.PID)
}
