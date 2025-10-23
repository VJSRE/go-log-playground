package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	levels  = []string{"INFO", "DEBUG", "ERROR", "WARNING", "CRITICAL", "ALERT", "EMERGENCY"}
	modules = []string{"AuthService", "PaymentGateway", "UserService", "InventoryManager", "APIHandler"}
	users   = []string{"alice", "bob", "charlie", "diana", "eve", "frank", "grace", "henry", "irene", "jack"}
	status  = []int{200, 201, 400, 401, 403, 404, 429, 500, 503}
)

func radmonIp(r *rand.Rand) string {
	return fmt.Sprintf("%d.%d.%d.%d", r.Intn(256), r.Intn(256), r.Intn(256), r.Intn(256))
}

func randomTimestamp(r *rand.Rand) string {
	now := time.Now()
	secBack := r.Int63n(30 * 24 * 60 * 60)
	ts := now.Add(-time.Duration(secBack) * time.Second)
	return ts.Format("2006-01-02 15:04:05,000")
}

func generateLogLine(r *rand.Rand) string {
	level := levels[r.Intn(len(levels))]
	module := modules[r.Intn(len(modules))]
	user := users[r.Intn(len(users))]
	ip := radmonIp(r)
	code := status[r.Intn(len(status))]

	msgIdx := r.Intn(6)
	var message string
	switch msgIdx {
	case 0:
		message = fmt.Sprintf("User %s logged in successfully from %s", user, ip)
	case 1:
		message = fmt.Sprintf("User %s failed authentication from %s", user, ip)
	case 2:
		message = fmt.Sprintf("Payment request processed for %s with status %d", user, code)
	case 3:
		message = fmt.Sprintf("Data sync initiated by %s on %s", user, module)
	case 4:
		message = fmt.Sprintf("API call to /v1/%s returned %d", strings.ToLower(module), code)
	default:
		message = fmt.Sprintf("System check by %s completed successfully", user)
	}

	levelPaded := fmt.Sprintf("%-8s", level)
	return fmt.Sprintf("%s %s [%s] %s\n ", randomTimestamp(r), levelPaded, module, message)

}

func main() {
	var lines int
	var outfile string
	var delay float64

	flag.IntVar(&lines, "lines", 10000, "Number of lines to generate")
	flag.StringVar(&outfile, "outfile", "synthetic_app.log", "File to write generated logs to")
	flag.Float64Var(&delay, "delay", 0.0, "Optional delay in seconds between synthetic operations")
	flag.Parse()

	f, err := os.Create(outfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create file %s: %s\n", outfile, err)
		os.Exit(1)
	}
	defer f.Close()

	w := bufio.NewWriterSize(f, 1<<20)
	defer w.Flush()

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	fmt.Printf("Generating %d log lines in %q...\n", lines, outfile)
	sleepDur := time.Duration(delay * float64(time.Second))

	for i := 0; i < lines; i++ {
		if _, err := w.WriteString(generateLogLine(r) + "\n"); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing log at line %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if sleepDur > 0 {
			time.Sleep(sleepDur)
		}
		if (i+1)%10000 == 0 {
			fmt.Printf("Processed %d lines in %q\n", i+1, outfile)
		}
	}

	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "flush error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Log generation completd !!")
}
