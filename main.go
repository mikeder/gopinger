package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	httpstat "github.com/tcnksm/go-httpstat"
)

// Check definition of a check to be performed by the program.
type Check struct {
	URL         string        `json:"url"`
	HealthyCode int           `json:"healthycode"`
	Method      string        `json:"method"`
	Timeout     time.Duration `json:"timeout"`
}

// Result is the output of a check against the given Check.
type Result struct {
	Code   int `json:"code"`
	Stats  httpstat.Result
	Status string `json:"status"`
	Reason string `json:"reason"`
	URL    string `json:"url"`
}

func main() {
	// TODO: Move these into db or something
	sites := []string{
		"https://forbar.net",
		"https://mikeder.net",
		"https://docker.sqweeb.net",
		"https://music.sqweeb.net",
		"https://git.sqweeb.net",
		"https://api.github.com",
	}
	// Setup list of checks to be performed
	var checks []Check
	var results []Result
	for _, url := range sites {
		checks = append(checks, Check{URL: url, HealthyCode: 200, Method: "GET", Timeout: 3})
	}
	go runChecks(&checks, &results)

	// Setup server to do things
	http.HandleFunc("/checks/run", func(w http.ResponseWriter, r *http.Request) {
		runChecks(&checks, &results)
		discardOldResults(&results)
	})

	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(results)
		if err != nil {
			fmt.Fprintf(w, "%v\n", err.Error())
		}
		fmt.Fprintf(w, "{\"results\":%v}\n", string(b))
	})

	http.ListenAndServe(":3001", nil)
}

func runChecks(c *[]Check, r *[]Result) {
	var client http.Client
	// Perform the checks
	for _, check := range *c {
		fmt.Println("Calling: " + check.URL)
		result, err := performCheck(&client, &check)
		if err != nil {
			fmt.Printf("%v \n", err.Error())
		}
		*r = append(*r, result)
	}
}

func discardOldResults(r *[]Result) {
	if len(*r) > 10 {
		var newR []Result
		newR = *r
		*r = newR[:10]
	}
}

func performCheck(cl *http.Client, ch *Check) (Result, error) {
	var result Result
	result.URL = ch.URL

	// Set custom timeout per check
	cl.Timeout = time.Duration(time.Second * ch.Timeout)
	req, err := http.NewRequest(ch.Method, ch.URL, nil)

	// Create a httpstat powered context
	var stats httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &stats)
	req = req.WithContext(ctx)

	// Perform the check and defer body close
	resp, err := cl.Do(req)
	if err != nil {
		result.Status = "FAIL"
		result.Reason = err.Error()
		return result, err
	}
	defer resp.Body.Close()

	result.Stats = stats
	result.Status = "PASS"
	result.Code = resp.StatusCode
	if result.Code != ch.HealthyCode {
		result.Status = "FAIL"
		result.Reason = "StatusCode Mismatch!"
	}
	return result, nil
}
