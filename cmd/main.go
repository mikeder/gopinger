package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mikeder/gopinger/internal/checker"
	httpstat "github.com/tcnksm/go-httpstat"
)

const version = "0.1"

// Check definition of a check to be performed by the program.
type Check struct {
	URL         string        `json:"url"`
	HealthyCode int           `json:"healthycode"`
	Method      string        `json:"method"`
	Timeout     time.Duration `json:"timeout"`
}

// Result is the output of a check against the given Check.
type Result struct {
	Code int `json:"code"`
	Stats
	Status string    `json:"status"`
	Reason string    `json:"reason"`
	URL    string    `json:"url"`
	Time   time.Time `json:"time"`
}

// Stats are http request duration metrics
type Stats struct {
	Lookup  time.Duration `json:"lookup"`
	Connect time.Duration `json:"connect"`
	Process time.Duration `json:"process"`
}

func main() {
	printVersion := flag.Bool("version", false, "print version for this thing")
	flag.Parse()
	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}
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
	results := make(map[string][]Result)
	for _, url := range sites {
		checks = append(checks, Check{URL: url, HealthyCode: 200, Method: "GET", Timeout: 3})
	}
	go runChecks(&checks, results)

	mc := checker.NewMinecraft("192.168.2.211", "25565", 10)
	mc.PerformCheck()

	// Setup server to do things
	http.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(*r)
		if err != nil {
			fmt.Println(w, err.Error())
		}
		fmt.Fprintln(w, string(b))
	})

	http.HandleFunc("/checks/run", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(checks)
		if err != nil {
			fmt.Fprintf(w, "%v\n", err.Error())
		}
		runChecks(&checks, results)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "{\"checks\":%v}\n", string(b))
		discardOldResults(results)
	})

	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(results)
		if err != nil {
			fmt.Fprintf(w, "%v\n", err.Error())
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "{\"results\":%v}\n", string(b))
	})
	http.ListenAndServe(":3001", nil)
}

func runChecks(c *[]Check, r map[string][]Result) {
	var client http.Client
	// Perform the checks
	for _, check := range *c {
		result, err := performCheck(&client, &check)
		if err != nil {
			fmt.Printf("%v \n", err.Error())
		}
		r[check.URL] = append(r[check.URL], result)
	}
}

func discardOldResults(r map[string][]Result) {
	for k, v := range r {
		if len(v) > 10 {
			r[k] = v[len(v)-10:]
		}
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
		result.Status = "fail"
		result.Reason = err.Error()
		return result, err
	}
	defer resp.Body.Close()

	result.Stats = Stats{
		Lookup:  stats.NameLookup / time.Millisecond,
		Connect: stats.Connect / time.Millisecond,
		Process: stats.ServerProcessing / time.Millisecond,
	}
	result.Time = time.Now()
	result.Status = "pass"
	result.Code = resp.StatusCode
	if result.Code != ch.HealthyCode {
		result.Status = "fail"
		result.Reason = "status code mismatch"
	}
	return result, nil
}
