package main

import (
	"flag"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/OCCASS/avito-intern/internal/config"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	cfgPath := flag.String("c", "", "Configuration file path.")
	flag.Parse()

	cfg := config.MustLoad(*cfgPath)

	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 30 * time.Second

	targeter := createTarget(cfg)

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Team add load test") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Println("Load test results:")
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Max latency: %s\n", metrics.Latencies.Max)
	fmt.Printf("Success rate: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Total requests: %d\n", metrics.Requests)
	fmt.Printf("Throughput: %.2f requests/sec\n", metrics.Throughput)
}

func createTarget(cfg *config.Config) vegeta.Targeter {
	var counter uint64
	startTime := time.Now().Unix()

	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		id := atomic.AddUint64(&counter, 1)
		name := fmt.Sprintf("team_%d_%d", startTime, id)
		u1IdWithTime := fmt.Sprintf("u_%d_%d", startTime, id*3+1)
		u2IdWithTime := fmt.Sprintf("u_%d_%d", startTime, id*3+2)
		u3IdWithTime := fmt.Sprintf("u_%d_%d", startTime, id*3+3)

		tgt.Method = "POST"
		tgt.URL = fmt.Sprintf(`http://%s/team/add`, cfg.Server.Address())
		tgt.Header = map[string][]string{
			"Content-Type": {"application/json"},
			"Accept":       {"application/json"},
		}

		body := fmt.Sprintf(`{
			"team_name": "%s",
			"members": [
				{"user_id": "%s", "username": "Ivan", "is_active": true},
				{"user_id": "%s", "username": "Vasya", "is_active": false},
				{"user_id": "%s", "username": "Lena", "is_active": false}
			]
		}`, name, u1IdWithTime, u2IdWithTime, u3IdWithTime)

		tgt.Body = []byte(body)

		return nil
	}
}
