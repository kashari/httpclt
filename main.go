package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type headerFlags []string

func (h *headerFlags) String() string {
	return fmt.Sprintf("%v", *h)
}

func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func makeRequest(client *http.Client, method, urlStr, body string, headers headerFlags, wg *sync.WaitGroup, requestNum int) {
	if wg != nil {
		defer wg.Done()
	}

	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(strings.ToUpper(method), urlStr, reqBody)
	if err != nil {
		log.Printf("[Request %d] Error creating request: %v\n", requestNum, err)
		return
	}

	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			log.Printf("[Request %d] Warning: Ignoring invalid header format: %q.", requestNum, h)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Request %d] Error sending request: %v\n", requestNum, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("[Request %d] Status: %s\n", requestNum, resp.Status)

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		log.Printf("[Request %d] Error reading response body: %v\n", requestNum, err)
	}
}

func main() {
	url := flag.String("url", "", "The URL for the HTTP request. (Required)")
	method := flag.String("method", "GET", "The HTTP method to use (e.g., GET, POST, PUT, DELETE).")
	body := flag.String("body", "", "The request body for POST, PUT, or PATCH requests.")
	requests := flag.Int("requests", 1, "The total number of requests to send.")
	perSecond := flag.Int("per-second", 0, "The maximum number of requests per second. 0 means no limit.")

	var headers headerFlags
	flag.Var(&headers, "header", "A request header in 'Key: Value' format. Can be specified multiple times.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A simple Go utility to make configurable HTTP requests.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s --url https://api.example.com/items --method POST --header \"Content-Type: application/json\" --body '{\"name\":\"new item\"}'\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nConcurrency Example:\n")
		fmt.Fprintf(os.Stderr, "  %s --url https://api.example.com/items --requests 100 --per-second 10\n", os.Args[0])
	}

	flag.Parse()

	if *url == "" {
		fmt.Fprintln(os.Stderr, "Error: --url flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	client := &http.Client{}

	if *requests <= 1 {
		var reqBody io.Reader
		if *body != "" {
			reqBody = strings.NewReader(*body)
		}

		req, err := http.NewRequest(strings.ToUpper(*method), *url, reqBody)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		for _, h := range headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) != 2 {
				log.Printf("Warning: Ignoring invalid header format: %q. Expected 'Key: Value'.", h)
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			req.Header.Set(key, value)
		}

		fmt.Printf("--> %s %s\n", req.Method, req.URL)
		for key, values := range req.Header {
			for _, value := range values {
				fmt.Printf("--> %s: %s\n", key, value)
			}
		}
		if *body != "" {
			fmt.Printf("--> Body: %s\n", *body)
		}
		fmt.Println("-->")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		fmt.Printf("<-- %s\n", resp.Status)

		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("<-- %s: %s\n", key, value)
			}
		}
		fmt.Println("<--")

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		fmt.Println(string(respBody))
		return 
	}

	var wg sync.WaitGroup
	var ticker <-chan time.Time

	if *perSecond > 0 {
		ticker = time.NewTicker(time.Second / time.Duration(*perSecond)).C
	}

	fmt.Printf("Sending %d requests to %s...\n", *requests, *url)
	startTime := time.Now()

	for i := 1; i <= *requests; i++ {
		if ticker != nil {
			<-ticker 
		}
		wg.Add(1)
		go makeRequest(client, *method, *url, *body, headers, &wg, i)
	}

	wg.Wait() 

	duration := time.Since(startTime)
	fmt.Printf("\nFinished %d requests in %v\n", *requests, duration)
}


