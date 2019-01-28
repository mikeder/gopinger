package main

import (
  "flag"
  "fmt"
  "net/http"
  "time"
  "github.com/mikeder/gopinger/lib"
  "github.com/jasonlvhit/gocron"
)

func main() {
    url := flag.String("url", "https://golang.org", "URL to perform check against.")
    code := flag.Int("code", 200, "Healthy response code to check for.")
    method := flag.String("method", "GET", "HTTP method to use in request.")
    debug := flag.Bool("v", false, "Print debug information.")
    flag.Parse()

    if *debug {
      fmt.Println("## Debug Info ##\n")
      fmt.Printf("URL: %v\n", *url)
      fmt.Printf("Code: %v\n", *code)
      fmt.Printf("Method: %v\n", *method)
    }

    client := &http.Client{
    	Timeout: time.Second * 5,
    }

    switch *method {
    case "GET":
      resp, err := client.Get(*url)
      if err != nil {
        panic(err)
      }
      handleResponse(*resp, *code, *url)
    default:
      panic("Requested method is not yet supported.")
    }
}

func performCheck(client http.Client, check checks.Check) checks.Result {
  result := checks.Result{
    Status: "good",
    Reason: "things worked!",
  }
  return result
}

func handleResponse(resp http.Response, code int, url string) {
  if resp.StatusCode != code {
    fmt.Printf("Health check failed for %v\n", url)
    fmt.Println("\n")
    fmt.Printf("Response: %v", resp)
  } else {
    fmt.Printf("Health check passed for %v\n", url)
    fmt.Println("\n")
    fmt.Printf("Response: %v", resp)
  }
}
