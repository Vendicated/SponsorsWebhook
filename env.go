package main

import (
	"fmt"
	"os"
	"strconv"
)

var (
	Secret     []byte
	Port       string
	WebhookUrl string
)

func init() {
	Secret = []byte(os.Getenv("SPONSORS_SECRET"))
	if len(Secret) == 0 {
		fmt.Println("Missing SPONSORS_SECRET env var")
		os.Exit(1)
	}

	WebhookUrl = os.Getenv("SPONSORS_WEBHOOK_URL")
	if WebhookUrl == "" {
		fmt.Println("Missing SPONSORS_WEBHOOK_URL env var")
		os.Exit(1)
	}

	Port = os.Getenv("SPONSORS_PORT")
	if Port == "" {
		Port = "1928"
		fmt.Println("Using port", Port+". Set SPONSORS_PORT env var to override")
	} else {
		if _, err := strconv.Atoi(Port); err != nil {
			fmt.Println("Invalid SPONSORS_PORT env var:", Port)
			os.Exit(1)
		}
	}
}
