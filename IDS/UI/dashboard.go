package main

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"
)

type DashboardData struct {
	CPU        string
	Memory     string
	Uptime     string
	GoVersion  string
	NumGoroute int
}

var startTime = time.Now()

func handler(w http.ResponseWriter, r *http.Request) {
	data := DashboardData{
		CPU:        "N/A", // Go does not have direct CPU usage tracking
		Memory:     fmt.Sprintf("%v MB", getMemoryUsage()),
		Uptime:     fmt.Sprintf("%s", time.Since(startTime).Round(time.Second)),
		GoVersion:  runtime.Version(),
		NumGoroute: runtime.NumGoroutine(),
	}

	tmpl, err := template.New("dashboard").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template Parsing Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func getMemoryUsage() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Alloc / 1024 / 1024 // Convert bytes to MB
}

var htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Go Dashboard</title>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
	<style>
		body { padding: 20px; font-family: Arial, sans-serif; }
		.card { margin: 10px; }
	</style>
</head>
<body>
	<div class="container">
		<h1 class="text-center">Go Dashboard</h1>
		<div class="row">
			<div class="col-md-4">
				<div class="card text-white bg-primary">
					<div class="card-body">
						<h5 class="card-title">CPU Usage</h5>
						<p class="card-text">{{.CPU}}</p>
					</div>
				</div>
			</div>
			<div class="col-md-4">
				<div class="card text-white bg-success">
					<div class="card-body">
						<h5 class="card-title">Memory Usage</h5>
						<p class="card-text">{{.Memory}}</p>
					</div>
				</div>
			</div>
			<div class="col-md-4">
				<div class="card text-white bg-warning">
					<div class="card-body">
						<h5 class="card-title">Uptime</h5>
						<p class="card-text">{{.Uptime}}</p>
					</div>
				</div>
			</div>
		</div>
		<div class="row">
			<div class="col-md-6">
				<div class="card text-white bg-info">
					<div class="card-body">
						<h5 class="card-title">Go Version</h5>
						<p class="card-text">{{.GoVersion}}</p>
					</div>
				</div>
			</div>
			<div class="col-md-6">
				<div class="card text-white bg-danger">
					<div class="card-body">
						<h5 class="card-title">Goroutines</h5>
						<p class="card-text">{{.NumGoroute}}</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</body>
</html>
`

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Dashboard is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
