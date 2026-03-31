package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"mimealogic/pkg"
)

func startWebInterface(engine *pkg.Engine, metrics *Metrics, stateManager *StateManager, config Config) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Explicitly tell the browser to use UTF-8
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Calculate live data for the dashboard
		lastWatered := stateManager.GetOrDefault(time.Now().Add(-24 * time.Hour))
		hoursPassed := time.Since(lastWatered).Hours()
		currentMoisture := engine.PredictMoisture(config.InitialMoisture, config.EvaporationRate, hoursPassed)

		// Start HTML output using entities for emojis
		// &#127807; = 🌿
		fmt.Fprintf(w, `
			<html>
			<head>
				<title>MimeaLogic Dashboard</title>
				<style>
					body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f7f6; color: #333; text-align: center; padding: 50px; }
					.card { background: white; border-radius: 15px; box-shadow: 0 4px 8px rgba(0,0,0,0.1); display: inline-block; padding: 30px; min-width: 350px; }
					.moisture-value { font-size: 54px; font-weight: bold; color: #2ecc71; margin: 10px 0; }
					.status-ok { color: #27ae60; font-weight: bold; font-size: 1.2em; }
					.status-warn { color: #e67e22; font-weight: bold; font-size: 1.2em; }
					.metric-box { margin-top: 25px; border-top: 2px solid #f4f7f6; padding-top: 20px; font-size: 0.95em; color: #555; text-align: left; }
					h1 { margin-top: 0; color: #2c3e50; }
				</style>
				<meta http-equiv="refresh" content="5"> 
			</head>
			<body>
				<div class="card">
					<h1>&#127807; MimeaLogic Live</h1>
					<p style="color: #7f8c8d; margin-bottom: 0;">Current Soil Moisture</p>
					<div class="moisture-value">%.1f%%</div>
		`, currentMoisture)

		// &#9888; = ⚠️ | &#9989; = ✅
		if currentMoisture < config.MoistureThreshold {
			fmt.Fprintf(w, "<p class='status-warn'>&#9888; Status: Watering Required</p>")
		} else {
			fmt.Fprintf(w, "<p class='status-ok'>&#9989; Status: Optimal</p>")
		}

		fmt.Fprintf(w, `
					<div class="metric-box">
						<p><b>Last Watered:</b> %s ago</p>
						<p><b>System Threshold:</b> %.1f%%</p>
						<p><b>Engine Status:</b> Nominal</p>
					</div>
				</div>
			</body>
			</html>
		`, formatDuration(hoursPassed), config.MoistureThreshold)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🌐 Web interface starting on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
