package jmsnet

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
)

//integrates godog and go test command
var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

var timeoutMillis int

var timingResult float64

var fetchURL = "http://www.cnn.com"

var receiptChannel chan float64

func iSupplyBadRequestCode() error {
	timingResult = BadRequest
	return nil
}

func badRequestIsReported() error {
	got := CFResult(fetchURL, timingResult)
	want := fetchURL + " " + "Bad Request"
	if got != want {
		return fmt.Errorf("Wanted %s but got %s", want, got)
	}
	return nil
}

func iSupplyTimeouterrorCode() error {
	timingResult = RequestTimeoutOrError
	return nil
}

func timeouterrorIsReported() error {
	got := CFResult(fetchURL, timingResult)
	want := fetchURL + " " + "Request Timeout or Error"
	if got != want {
		return fmt.Errorf("Wanted %s but got %s", want, got)
	}
	return nil
}

func iSupplyValidTime(arg1 int) error {
	timingResult = float64(arg1)
	return nil
}

func fetchTimeIsReported() error {
	got := CFResult(fetchURL, timingResult)
	want := fetchURL

	want += " responded in "
	want += fmt.Sprintf("%f", timingResult)
	want += " seconds"
	if got != want {
		return fmt.Errorf("Wanted %s but got %s", want, got)
	}
	return nil
}

func iSpecifyAValidURLLike(arg1 string) error {
	fetchURL = arg1
	return nil
}

func iSpecifyATimeoutOfMs(arg1 int) error {
	timeoutMillis = arg1
	return nil
}

func iMeasureResponseTime() error {
	receiptChannel = make(chan float64)
	go ResponseTime(fetchURL, receiptChannel, timeoutMillis)
	for {
		select {
		case timingResult = <-receiptChannel:
			return nil
		case <-time.After(time.Duration(timeoutMillis) * time.Millisecond):
			return fmt.Errorf("Fetch from URL %s timed out in %d ms", fetchURL, timeoutMillis)
		}
	}
}

func positiveFetchTimeLessThanTimeoutIsReported() error {
	timeoutSeconds := float64(timeoutMillis / 1000)
	if timingResult <= float64(0) {
		return fmt.Errorf("receipt time of %f s must be greater than 0", timingResult)
	}
	if timingResult > timeoutSeconds {
		return fmt.Errorf("receipt time of %f s must be less than timeout of %f s", timingResult, timeoutSeconds)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I supply bad request code$`, iSupplyBadRequestCode)
	s.Step(`^bad request is reported$`, badRequestIsReported)
	s.Step(`^I supply timeouterror code$`, iSupplyTimeouterrorCode)
	s.Step(`^timeouterror is reported$`, timeouterrorIsReported)
	s.Step(`^I supply valid time (\d+)$`, iSupplyValidTime)
	s.Step(`^fetch time is reported$`, fetchTimeIsReported)
	s.Step(`^I specify a valid URL like "([^"]*)"$`, iSpecifyAValidURLLike)
	s.Step(`^I specify a timeout of (\d+) ms$`, iSpecifyATimeoutOfMs)
	s.Step(`^I measure response time$`, iMeasureResponseTime)
	s.Step(`^positive fetch time less than timeout is reported$`, positiveFetchTimeLessThanTimeoutIsReported)
}

//integrates godog and go test command
func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
