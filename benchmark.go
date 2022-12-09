package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/akamensky/argparse"
)

type CreateSessionRequestData struct {
	Username string `json:"username"`
}

type SendCodeRequestData struct {
	Username string   `json:"username"`
	Code     []string `json:"code"`
}

func start_session(url *string, id int, c chan float64) {
	username := fmt.Sprintf("user-%d", id)
	user := CreateSessionRequestData{
		Username: username,
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	request, error := http.NewRequest("POST", *url+"/create_session", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	start := time.Now()
	response, error := client.Do(request)
	elapsed := time.Since(start)

	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Printf("user - %d\tresponse(/start_session) - %d\telapsed - %s\n", id, response.StatusCode, elapsed)

	if response.StatusCode != 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			panic("Create session failed!")
		}
		fmt.Println("Request body:", string(body))

		// Print the request headers
		fmt.Println("Request headers:")
		for key, value := range response.Header {
			fmt.Println(key+":", value)
		}

		panic("Create session failed!")
	}

	c <- elapsed.Seconds()
}

func send_code(url *string, id int, runs int, c chan float64) {
	username := fmt.Sprintf("user-%d", id)
	code := SendCodeRequestData{
		Username: username,
		Code:     []string{"move(Direction.RIGHT)"},
	}

	jsonData, err := json.Marshal(code)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	var average_time float64 = 0

	for i := 0; i < int(runs); i++ {
		request, err := http.NewRequest("POST", *url+"/send_code", bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		start := time.Now()
		response, err := client.Do(request)
		elapsed := time.Since(start)

		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		fmt.Printf("user - %d\tresponse(/send_code) - %d\telapsed - %s\n", id, response.StatusCode, elapsed)

		if response.StatusCode != 200 {
			// Print the request body
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				panic("Send code failed!")
			}
			fmt.Println("Request body:", string(body))

			// Print the request headers
			fmt.Println("Request headers:")
			for key, value := range response.Header {
				fmt.Println(key+":", value)
			}

			panic("Send code failed!")
		}

		average_time += elapsed.Seconds()
	}

	fmt.Printf("%f\n", (average_time / float64(runs)))
	c <- (average_time / float64(runs))
}

func main() {
	parser := argparse.NewParser("lynx-benchmark", "Benchmark mainly lynx-runner + lynx-runtime")

	url := parser.StringPositional(&argparse.Options{Default: "http://server.blazej-smorawski.com", Help: "Url that will be queried at `/start_session` and `/send_code`"})
	threads_ptr := parser.Int("t", "threads", &argparse.Options{Default: 8, Help: "Count of threads to be spawned which will send parallel queries to `runner`"})
	runs_ptr := parser.Int("r", "runs", &argparse.Options{Default: 8, Help: "Number of times each `send_code` query should be repeated by each thread"})
	only_code_ptr := parser.Flag("c", "only-code", &argparse.Options{Help: "Do not create new sessions, just use `send_code`"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	threads := *threads_ptr
	runs := *runs_ptr
	only_code := *only_code_ptr

	var c chan float64 = make(chan float64)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !only_code {
		for i := 0; i < threads; i++ {
			go start_session(url, i, c)
		}

		var create_session_avg float64 = 0
		for i := 0; i < threads; i++ {
			create_session_avg += <-c
		}
		create_session_avg /= float64(threads)

		fmt.Printf("Average session creation time: %fs\n", create_session_avg)
	}

	for i := 0; i < threads; i++ {
		go send_code(url, i, runs, c)
	}

	var send_code_avg float64 = 0
	for i := 0; i < threads; i++ {
		send_code_avg += <-c
	}
	send_code_avg /= float64(threads)

	fmt.Printf("Average code execution time: %fs\n", send_code_avg)
}
