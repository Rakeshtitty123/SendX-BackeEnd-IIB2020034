package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Cache is a struct that stores the HTML of a webpage and the timestamp
// when it was last requested
type Cache struct {
	HTML      string
	Timestamp time.Time
}

var cache = make(map[string]Cache)

// PriorityWorker is a struct that stores a URL, a retry limit, and whether the customer is paying
type PriorityWorker struct {
	URL            string
	RetryLimit     int
	CustomerPaying bool
}

var PayingWorkQueue = make(chan PriorityWorker)
var NonPayingWorkQueue = make(chan PriorityWorker)

// WorkerPool is a slice of Worker channels
var WorkerPool []chan PriorityWorker

// MaxWorkers is the maximum number of active workers
const MaxWorkers = 10

func main() {
	// Create the worker pool
	for i := 0; i < MaxWorkers; i++ {
		WorkerPool = append(WorkerPool, make(chan PriorityWorker))
	}

	// Create a cache folder if it doesn't exist to store downloaded files
	err := os.Mkdir("cache", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating cache folder: %v\n", err)
	}
	fmt.Println("Created cache folder")

	// Start the worker pipeline
	go workerPipeline()

	http.HandleFunc("/download", downloadHandler)
	http.ListenAndServe(":8080", nil)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the query string to get the URL, retry limit, and customer paying status
	query := r.URL.Query()
	url := query.Get("url")
	retryLimit, _ := strconv.Atoi(query.Get("retry_limit"))
	if retryLimit > 10 {
		retryLimit = 10
	}
	customerPaying, _ := strconv.ParseBool(query.Get("customer_paying"))

	// Check if the webpage has been requested in the last 24 hours
	if c, ok := cache[url]; ok && time.Since(c.Timestamp) < 24*time.Hour {
		// Serve the webpage from the cache
		c.Timestamp = time.Now()
		w.Write([]byte(c.HTML))
		return
	}

	// Create a PriorityWorker based on whether the customer is paying or not
	worker := PriorityWorker{
		URL:            url,
		RetryLimit:     retryLimit,
		CustomerPaying: customerPaying,
	}

	// Determine the appropriate work queue to send the worker to
	var workQueue chan PriorityWorker
	if customerPaying {
		workQueue = PayingWorkQueue
	} else {
		workQueue = NonPayingWorkQueue
	}

	// Send the worker to the selected work queue
	workQueue <- worker
}

func workerPipeline() {
	for {
		select {
		case worker := <-PayingWorkQueue:
			handleWorker(worker)
		case worker := <-NonPayingWorkQueue:
			handleWorker(worker)
		}
	}
}

func handleWorker(worker PriorityWorker) {
	// Get the next worker from the worker channel
	// Download the webpage
	start := time.Now()
	html, err := downloadWebpage(worker.URL, worker.RetryLimit)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Cache the webpage
	cache[worker.URL] = Cache{
		HTML:      html,
		Timestamp: time.Now(),
	}

	// Download the webpage to the local file system
	err = saveWebpage(html, worker.URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("Customer (Paying: %v) served in %s\n", worker.CustomerPaying, elapsed)
}

func downloadWebpage(url string, retryLimit int) (string, error) {
	retries := 0
	for {
		resp, err := http.Get(url)
		if err != nil {
			if retries < retryLimit {
				retries++
				continue
			}
			return "", fmt.Errorf("error downloading webpage: %v", err)
		}
		defer resp.Body.Close()

		html, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if retries < retryLimit {
				retries++
				continue
			}
			return "", fmt.Errorf("error reading webpage body: %v", err)
		}

		return string(html), nil
	}
}

func saveWebpage(html, url string) error {
	replacer := strings.NewReplacer("/", "_", ":", "^")
	updatedURL := replacer.Replace(url)
	file, err := os.Create("cache/" + fmt.Sprintf("%s.html", updatedURL))
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, html)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
