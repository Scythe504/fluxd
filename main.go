package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type WorkerError struct {
	err_code    string
	err_message string
}

func randInt() string {
	max := big.NewInt(255)
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}

	parsedInt := strconv.Itoa(int(r.Uint64()))

	return parsedInt
}

func main() {
	cmd := exec.Command("./test.sh")
	wg := sync.WaitGroup{}

	dlq := make(chan WorkerError, 1)
	fmt.Println("---Started---")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal("failed to created in pipe")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("failed to create out pipe")
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		log.Fatal(err.Error())
	}

	scanner := bufio.NewScanner(stdout)
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	fmt.Println("---Scanner Initiated---")

	fmt.Println("---Scanner started its work---")
	wg.Go(func() {
		for err := range dlq {
			fmt.Println(err)
		}
	})
	wg.Go(func() {
		defer close(dlq)
		for scanner.Scan() {
			txt := strings.Split(scanner.Text(), ":")
			if txt[0] == "ERR" {
				dlq <- WorkerError{
					err_code:    txt[1],
					err_message: txt[len(txt)-1],
				}
			}
			
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	})

	fmt.Println("---Random Int Fed to Test Script---")
	for range 16 {
		randCode := randInt()
		n, err := io.WriteString(stdin, randCode+"\n")
		if err != nil {
			fmt.Println("N: \n", n, "ERR: \n", err)
		}
	}
	stdin.Close()
	dlq <- WorkerError{
		err_code: "100",
		err_message: "ERROR_MANUAL",
	}

	fmt.Println("\n---Channel Closed---?")
	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
}
