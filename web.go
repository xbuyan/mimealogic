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
		// Calculate live data for the dashboard
		lastWatered := stateManager.GetOrDefault(time.Now().Add(-24 * time.Hour))
		hoursPassed := time.Since(lastWatered).Hours()
		currentMoisture := engine.PredictMoisture(config.InitialMoisture, config.EvaporationRate, hoursPassed)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<html>
			<head>
				<title>MimeaLogic Dashboard</title>
				<style>
					body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f7f6; color: #333; text-align: center; padding: 50px; }
					.card { background: white; border-radius: 15px; box-shadow: 0 4px 8px rgba(0,0,0,0.1); display: inline-block; padding: 30px; min-width: 300px; }
					.moisture-value { font-size: 48px; font-weight: bold; color: #2ecc71; }
					.status-ok { color: #27ae60; font-weight: bold; }
					.status-warn { color: #e67e22; font-weight: bold; }
					.metric-box { margin-top: 20px; border-top: 1px solid #eee; padding-top: 20px; font-size: 0.9em; color: #666; }
				</style>
				<meta http-equiv="refresh" content="5"> 
			</head>
			<body>
				<div class="card">
					<h1>🌿 MimeaLogic Live</h1>
					<p>Current Soil Moisture</p>
					<div class="moisture-value">%.1f%%</div>
		`, currentMoisture)

		if currentMoisture < config.MoistureThreshold {
			fmt.Fprintf(w, "<p class='status-warn'>⚠️ Status: Watering Required</p>")
		} else {
			fmt.Fprintf(w, "<p class='status-ok'>✅ Status: Optimal</p>")
		}

		// Add logic stats
		fmt.Fprintf(w, `
					<div class="metric-box">
						<p><b>Last Watered:</b> %s ago</p>
						<p><b>System Threshold:</b> %.1f%%</p>
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
	http.ListenAndServe(":"+port, nil)
}
