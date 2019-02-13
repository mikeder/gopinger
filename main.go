package main

import (
	"fmt"
	"log"
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
	Code     int     `json:"code"`
	Duration float32 `json:"duration"`
	Status   string  `json:"status"`
	Reason   string  `json:"reason"`
}

func main() {

	sites := []string{"https://mikeder.net", "https://sqweeb.net", "https://psymux.net"}

	// Setup list of checks to be performed
	var checks []Check
	for _, url := range sites {
		checks = append(checks, Check{URL: url, HealthyCode: 200, Method: "GET", Timeout: 3})
	}

	var results []Result
	var client http.Client
	// Perform the checks
	for _, check := range checks {
		fmt.Println("Calling: " + check.URL)
		results = append(results, performCheck(client, check))
	}

	for _, result := range results {
		fmt.Println(result.Status, result.Code, result.Duration)
	}
}

func performCheck(cl http.Client, ch Check) Result {

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
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Show the results
	log.Printf("DNS lookup: %d ms", int(stats.DNSLookup/time.Millisecond))
	log.Printf("TCP connection: %d ms", int(stats.TCPConnection/time.Millisecond))
	log.Printf("TLS handshake: %d ms", int(stats.TLSHandshake/time.Millisecond))
	log.Printf("Server processing: %d ms", int(stats.ServerProcessing/time.Millisecond))
	log.Printf("Content transfer: %d ms", int(stats.ContentTransfer(time.Now())/time.Millisecond))

	var result Result
	result.Code = resp.StatusCode
	result.Duration = float32((stats.DNSLookup + stats.TCPConnection + stats.TLSHandshake + stats.ServerProcessing + stats.Connect) * time.Millisecond)
	result.Status = "PASS"
	if result.Code != ch.HealthyCode {
		result.Status = "FAIL"
		result.Reason = "StatusCode Mismatch!"
	}

	return result
}
