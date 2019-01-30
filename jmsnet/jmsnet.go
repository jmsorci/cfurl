//Package jmsnet provides utility functions
package jmsnet

import (
	"fmt"
	"net/http"
	"time"
)

//BadRequest is a status code -1
const BadRequest float64 = float64(-1)

//RequestTimeoutOrError is a status code -2
const RequestTimeoutOrError float64 = float64(-2)

/*
ResponseTime performs a GET request against url with timeout (ms)
Sends the time in (s) back along the supplied channel c
*/
func ResponseTime(url string, c chan float64, to int) {

	client := &http.Client{
		Timeout: time.Duration(to) * time.Millisecond,
	}

	start := time.Now()

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		c <- BadRequest
	}

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
		c <- RequestTimeoutOrError
	}

	defer response.Body.Close()

	elapsed := time.Since(start).Seconds()

	c <- elapsed
}

//CFResult prints the timing result given url and time
func CFResult(url string, result float64) string {
	msg := url
	msg += " "
	switch result {
	case BadRequest:
		msg += "Bad Request"
	case RequestTimeoutOrError:
		msg += "Request Timeout or Error"
	default:
		msg += "responded in "
		msg += fmt.Sprintf("%f", result)
		msg += " seconds"
	}
	return msg
}
