package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/jmsorci/cfurl/jmsnet"
)

func flagUsage() {
	usageText := `cfurl is acli tool that compares response times (s) url1 versus url2 
	with an optionally-specified timeout (ms)
        
Usage:
cfurl -u1 url1  -u2 url2 [-t ms]

Example
cfurl -u1 http://www.cnn.com -u2 http://www.foxnews.com -t 2000`
	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func main() {

	flag.Usage = flagUsage

	u1 := flag.String("u1", "", "First URL to compare")
	u2 := flag.String("u2", "", "Second URL to compare")
	t := flag.Int("t", 2000, "Timeout in ms. Defaults to 2000 ms")

	flag.Parse()

	if *u1 == "" {
		fmt.Fprintf(os.Stderr, "-u1 required\n")
		flagUsage()
		return
	}

	if *u2 == "" {
		fmt.Fprintf(os.Stderr, "-u2 required\n")
		flagUsage()
		return
	}

	channel1 := make(chan float64)
	channel2 := make(chan float64)

	go jmsnet.ResponseTime(*u1, channel1, *t)
	go jmsnet.ResponseTime(*u2, channel2, *t)

	var result1 float64
	var result2 float64

	for {
		select {
		case result1 = <-channel1:
			fmt.Println(jmsnet.CFResult(*u1, result1))
		case result2 = <-channel2:
			fmt.Println(jmsnet.CFResult(*u2, result2))
		}

		if result1 != float64(0) && result2 != float64(0) {
			if result1 > float64(0) && result2 > float64(0) {
				fmt.Println("Difference in response time is: ", fmt.Sprintf("%f seconds", math.Abs(result1-result2)))
			}
			return
		}
	}

}
