package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	weatherStationURL = "http://windslammer.net/cgi-bin/ws.cgi"
	metricsPath       = "/metrics"
	defaultPort       = "9307"
)

var (
	metrics = map[string]prometheus.Gauge{
		"wind_dir":   prometheus.NewGauge(prometheus.GaugeOpts{Name: "windslammer_wind_direction_degrees", Help: "Current wind direction in degrees"}),
		"wind_speed": prometheus.NewGauge(prometheus.GaugeOpts{Name: "windslammer_wind_speed_mph", Help: "Current wind speed in miles per hour"}),
		"temp_lo":    prometheus.NewGauge(prometheus.GaugeOpts{Name: "windslammer_temperature_lower_fahrenheit", Help: "Temperature from lower elevation station in Fahrenheit"}),
		"temp_hi":    prometheus.NewGauge(prometheus.GaugeOpts{Name: "windslammer_temperature_upper_fahrenheit", Help: "Temperature from upper elevation station in Fahrenheit"}),
		"elev_lo":    prometheus.NewGauge(prometheus.GaugeOpts{Name: "windslammer_elevation_lower_feet", Help: "Elevation of lower weather station in feet"}),
		"elev_hi":    prometheus.NewGauge(prometheus.GaugeOpts{Name: "windslammer_elevation_upper_feet", Help: "Elevation of upper weather station in feet"}),
	}
)

func init() {
	for _, gauge := range metrics {
		prometheus.MustRegister(gauge)
	}
}

func fetchWeatherData() error {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Post(weatherStationURL, "application/x-www-form-urlencoded", strings.NewReader("SNAPSHOT"))
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("weather station returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	return parseAndUpdateMetrics(string(body))
}

func parseAndUpdateMetrics(data string) error {
	for _, part := range strings.Split(data, ",") {
		kv := strings.Split(strings.TrimSpace(part), "=")
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		valueStr := strings.TrimSpace(kv[1])

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			log.Printf("Failed to parse value for %s: %v", key, err)
			continue
		}

		if gauge, exists := metrics[key]; exists {
			gauge.Set(value)
		}
	}

	log.Printf("Updated metrics from windslammer data")
	return nil
}


func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":" + defaultPort
}

func main() {
	port := getPort()
	log.Println("Starting Windslammer Prometheus Exporter")

	http.HandleFunc(metricsPath, func(w http.ResponseWriter, r *http.Request) {
		if err := fetchWeatherData(); err != nil {
			log.Printf("Error fetching data for metrics: %v", err)
		}
		promhttp.Handler().ServeHTTP(w, r)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server starting on port %s", port)
	log.Printf("Metrics available at http://localhost%s%s", port, metricsPath)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
