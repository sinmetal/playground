package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/sinmetal/gcp_playground/bigtable"

	"cloud.google.com/go/trace"
)

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	fmt.Fprintf(w, "Hello, Google Cloud Platform Playground")
	end := time.Now()

	fmt.Printf("duration=%dns\n", end.Sub(start).Nanoseconds())
}

func main() {
	ctx := context.Background()

	project := flag.String("project", "", "The Google Cloud Platform project ID. Required.")
	bigtableInstance := flag.String("bigtableInstance", "", "The Google Cloud Bigtable instance ID. Required.")
	flag.Parse()

	for _, f := range []string{"project"} {
		if flag.Lookup(f).Value.String() == "" {
			log.Fatalf("The %s flag is required.", f)
		}
	}
	fmt.Printf("project=%s\n", *project)
	fmt.Printf("bigtableInstance=%s\n", *bigtableInstance)

	bigtable.SetUp(*project, *bigtableInstance)

	tc, err := trace.NewClient(ctx, *project)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	http.HandleFunc("/", handler)
	http.Handle("/bigtable", tc.HTTPHandler(http.HandlerFunc(bigtable.HandlerBigtable)))

	fmt.Println("listen start")
	http.ListenAndServe(":8080", nil)
}