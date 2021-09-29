package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	gethUrl := os.Getenv("GETH_URL")
	if gethUrl == "" {
		log.Fatal("GETH_URL environmental variable is required")
	}
	stopHeightString := os.Getenv("STOP_HEIGHT")
	if stopHeightString == "" {
		log.Fatal("STOP_HEIGHT environmental variable is required")
	}
	stopHeight, err := strconv.ParseUint(stopHeightString, 10, 64)
	if err != nil {
		log.Fatalf("STOP_HEIGHT isn't a number: %v\n", err)
	}
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		log.Fatal("NAMESPACE environmental variable is required")
	}

	client, err := ethclient.Dial(gethUrl)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum url: %s\n", err)
	}
	block, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Error getting the initial block height: %s", err)
	}
	fmt.Println(block)
	if block > stopHeight {
		log.Fatalf("The client has already passed %d\n", stopHeight)
	}
	for {
		block, err := client.BlockNumber(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(block)
		if stopHeight <= block {
			fmt.Println("Target height reached, stopping replica!")
			cmd := exec.Command("kubectl", "scale", "sts", "l2geth-replica", "--replicas=0", "-n", namespace)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Scaled down: %q\n", out.String())
			os.Exit(0)
		}
		time.Sleep(8 * time.Second)

	}
}
