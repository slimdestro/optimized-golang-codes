/**
multiple endpoints (via a config.yaml), each endpoint will have their own rates (because the providers have rate limiting)
and with staying in that rate limit it will do some http requests to the providers with the help of http tracer, 
will give us a report. - We will need to run this on multipler regions for 1 time. - Collect all data and write a good report (visualize the data using charts)

*/
package main

import (
 "fmt"
 "io/ioutil"
 "net/http"
 "time"
 "encoding/json"
 "github.com/gorilla/mux"
 "github.com/gorilla/handlers"
 "github.com/opentracing/opentracing-go"
 "github.com/uber/jaeger-client-go"
 "github.com/uber/jaeger-client-go/config"
 "github.com/uber/jaeger-lib/metrics"
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config holds the configuration for the application
type Config struct {
 Endpoints []Endpoint `yaml:"endpoints"`
 Regions []string `yaml:"regions"`
}

// Endpoint holds the configuration for each endpoint
type Endpoint struct {
 Name string `yaml:"name"`
 URL string `yaml:"url"`
 Rate int `yaml:"rate"`
 Region string `yaml:"region"`
}

// Report holds the data for the report
type Report struct {
 Name string `json:"name"`
 URL string `json:"url"`
 Region string `json:"region"`
 Data string `json:"data"`
}

// Reports holds the data for all the reports
type Reports []Report

// Data holds the data for each request
type Data struct {
 Name string `json:"name"`
 URL string `json:"url"`
 Region string `json:"region"`
 Data string `json:"data"`
}

// Data holds the data for all the requests
type Data []Data

// main is the entry point for the application
func main() {
 // Read the config file
 config := readConfig()

 // Setup the tracer
 tracer, closer := initJaeger("http-tracer")
 defer closer.Close()
 opentracing.SetGlobalTracer(tracer)

 // Setup the router
 router := mux.NewRouter()

 // Setup the prometheus metrics
 prometheus.MustRegister(prometheus.NewCounterFunc(prometheus.CounterOpts{
  Name: "http_requests_total",
  Help: "Total number of HTTP requests made",
 }, func() float64 {
  return float64(len(Data))
 }))

 // Setup the routes
 router.Handle("/metrics", promhttp.Handler())
 router.HandleFunc("/report", generateReport(config)).Methods("GET")

 // Start the server
 fmt.Println("Starting server on port 8080")
 http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router))
}

// readConfig reads the config.yaml file
func readConfig() Config {
 // Read the config file
 data, err := ioutil.ReadFile("config.yaml")
 if err != nil {
  panic(err)
 }

 // Unmarshal the config
 var config Config
 err = yaml.Unmarshal(data, &config)
 if err != nil {
  panic(err)
 }

 return config
}

// initJaeger initializes the Jaeger tracer
func initJaeger(service string) (opentracing.Tracer, io.Closer) {
 cfg := &config.Configuration{
  Sampler: &config.SamplerConfig{
   Type: "const",
   Param: 1,
  },
  Reporter: &config.ReporterConfig{
   LogSpans: true,
  },
 }
 tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
 if err != nil {
  panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
 }
 return tracer, closer
}

// generateReport generates the report
func generateReport(config Config) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  // Create a span for the request
  span := opentracing.StartSpan("generate_report")
  defer span.Finish()

  // Create a channel to receive the data
  dataChan := make(chan Data)

  // Iterate over the endpoints
  for _, endpoint := range config.Endpoints {
   // Iterate over the regions
   for _, region := range config.Regions {
    // Check if the endpoint is in the region
    if endpoint.Region == region {
     // Create a rate limiter
     limiter := time.Tick(time.Duration(endpoint.Rate) * time.Second)

     // Make the request
     go func() {
      // Wait for the rate limiter
      <-limiter

      // Create a span for the request
      span := opentracing.StartSpan("http_request", opentracing.ChildOf(span.Context()))
      defer span.Finish()

      // Make the request
      resp, err := http.Get(endpoint.URL)
      if err != nil {
       panic(err)
      }
      defer resp.Body.Close()

      // Read the response body
      body, err := ioutil.ReadAll(resp.Body)
      if err != nil {
       panic(err)
      }

      // Send the data to the channel
      dataChan <- Data{
       Name: endpoint.Name,
       URL: endpoint.URL,
       Region: region,
       Data: string(body),
      }
     }()
    }
   }
  }

  // Create a slice to hold the data
  var data Data

  // Receive the data from the channel
  for {
   select {
   case d := <-dataChan:
    data = append(data, d)
    if len(data) == len(config.Endpoints)*len(config.Regions) {
     break
    }
   }
  }

  // Create the report
  report := Reports{}
  for _, d := range data {
   report = append(report, Report{
    Name: d.Name,
    URL: d.URL,
    Region: d.Region,
    Data: d.Data,
   })
  }

  // Marshal the report
  reportJSON, err := json.Marshal(report)
  if err != nil {
   panic(err)
  }

  // Write the response
  w.Header().Set("Content-Type", "application/json")
  w.Write(reportJSON)
 }
}