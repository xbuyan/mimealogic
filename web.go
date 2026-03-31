package main

import (
	"fmt"
	"net/http"
	"os"
)

func startWebInterface(metrics *Metrics) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body style='font-family:sans-serif; text-align:center; padding-top:50px;'>")
		fmt.Fprintf(w, "<h1>🌿 MimeaLogic Agent: ONLINE</h1>")
		fmt.Fprintf(w, "<p>The plant logic engine is currently running in the background.</p>")
		fmt.Fprintf(w, "<p><b>System Status:</b> Operational</p>")
		fmt.Fprintf(w, "</body></html>")
	})

	port := os.Getenv("PORT") // Railway provides the port via environment variable
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🌐 Web interface starting on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
