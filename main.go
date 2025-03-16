package main

import (
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Get the list of hosts from the PINGER_DOMAINS environment variable
func getHostsFromEnv() []string {
    envHosts := os.Getenv("PINGER_DOMAINS")
    if envHosts == "" {
        fmt.Println("PINGER_DOMAINS is not set, using default values")
        return []string{"vk.com", "google.com", "yandex.ru", "api.telegram.org"}
    }
    return strings.Split(envHosts, ",")
}

// Registering metrics
var (
    pingDuration = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "ping_duration_seconds",
            Help: "Ping response time in seconds",
        },
        []string{"host"},
    )
)

// Function to perform ping
func ping(host string) {
    for {
        // Execute ping command
        out, err := exec.Command("ping", "-c", "1", host).Output()
        if err != nil {
            fmt.Printf("Error pinging %s: %v\n", host, err)
            pingDuration.WithLabelValues(host).Set(0)
        } else {
            // Extract response time from ping output
            re := regexp.MustCompile(`time=(\d+\.\d+)`)
            match := re.FindStringSubmatch(string(out))
            if len(match) > 1 {
                ms, err := strconv.ParseFloat(match[1], 64)
                if err == nil {
                    // Convert milliseconds to seconds and update metric
                    pingDuration.WithLabelValues(host).Set(ms / 1000)
                    fmt.Printf("Ping to %s: %f ms\n", host, ms)
                }
            }
        }
        // Wait before the next ping
        time.Sleep(10 * time.Second)
    }
}

func main() {
    // Get hosts from environment variable
    hosts := getHostsFromEnv()

    // Start pinging each host in a separate goroutine
    for _, host := range hosts {
        go ping(host)
    }

    // Expose metrics at /metrics endpoint
    http.Handle("/metrics", promhttp.Handler())
    fmt.Println("Metrics server running at :2112/metrics")
    http.ListenAndServe(":2112", nil)
}

