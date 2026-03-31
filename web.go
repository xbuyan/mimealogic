package main

import (
	"fmt"
	"net/http"
	"os"
)

func startWebInterface(metrics *Metrics) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>🌿 MimeaLogic Live Status</h1>")
		fmt.Fprintf(w, "<p>Agent is currently: <b>RUNNING</b></p>")
		// Add more stats here if you have time!
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🌐 Web interface starting on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
