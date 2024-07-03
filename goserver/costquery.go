package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
//	"sync"
	"time"
//	"strconv"

)

func makeHTTPSRequest(url string) ([]byte, error) {
	client := &http.Client{}
	maxRetries := 5
	retryDelay := 200 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		resp, err := client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				return body, nil
			}
		}

		fmt.Printf("Attempt %d failed: %v\n", i+1, err)

		if i < maxRetries-1 {
			time.Sleep(retryDelay)
			retryDelay += 200 * time.Millisecond
		}
	}

	return nil, fmt.Errorf("failed after %d attempts", maxRetries)
}

