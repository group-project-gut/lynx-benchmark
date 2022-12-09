package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type CreateSessionRequestData struct {
	Username string `json:"username"`
}

type SendCodeRequestData struct {
	Username string   `json:"username"`
	Code     []string `json:"code"`
}

func start_session(url *string, id uint64, c chan float64) {
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
		// Print the request body
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

func send_code(url *string, id uint64, runs uint64, c chan float64) {
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

	var average_time float64 = 0
	client := &http.Client{}
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
	args := os.Args[1:]

	var c chan float64 = make(chan float64)
	var url string = args[0]

	thread_count, err := strconv.ParseUint(args[1], 10, 64)
	runs_count, err := strconv.ParseUint(args[2], 10, 64)

	if err != nil {
		fmt.Println(err)
		return
	}

	for i := uint64(0); i < thread_count; i++ {
		go start_session(&url, i, c)
		time.Sleep(1 * time.Second)
	}

	var create_session_avg float64 = 0
	for i := uint64(0); i < thread_count; i++ {
		create_session_avg += <-c
	}
	create_session_avg /= float64(thread_count)

	fmt.Printf("Average session creation time: %fs\n", create_session_avg)

	for i := uint64(0); i < thread_count; i++ {
		go send_code(&url, i, runs_count, c)
	}

	var send_code_avg float64 = 0
	for i := uint64(0); i < thread_count; i++ {
		send_code_avg += <-c
	}
	send_code_avg /= float64(thread_count)

	fmt.Printf("Average code execution time: %fs\n", send_code_avg)
}
