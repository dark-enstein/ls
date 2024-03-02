package latency

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"time"
)

var (
	Servers = []string{
		"0.africa.pool.ntp.org",
		"0.asia.pool.ntp.org",
		"0.europe.pool.ntp.org",
		"0.north-america.pool.ntp.org",
		"0.oceania.pool.ntp.org",
		"0.south-america.pool.ntp.org",
		"pool.ntp.org",
	}
	expectedPingOutput = " 49.475/55.499/66.155/6.805 ms"
)

type LatencyFunc func(url string) (time.Duration, error)

var (
	PingCmd LatencyFunc = func(url string) (time.Duration, error) {
		out, err := exec.Command("ping", url, "-c 5", "-i 3", "|", "cut", "-d=", "-f 2").Output()
		if err != nil {
			return -1, err
		}
		// " 49.475/55.499/66.155/6.805 ms"
		// full: round-trip min/avg/max/stddev = 49.475/55.499/66.155/6.805 ms

		if bytes.Contains(out, []byte("Destination Host Unreachable")) {
			return -1, fmt.Errorf("destination host unreachable")
		}

		out = bytes.TrimSpace(out)
		out, after, f := bytes.Cut(out, []byte(" "))
		if !f || len(after) == 0 {
			return -1, fmt.Errorf("returned output did not match format: %s\n", expectedPingOutput)
		}
		oout := bytes.Split(out, []byte("/"))
		if len(oout) != 4 {
			return -1, fmt.Errorf("returned output did not match format: %s\n", expectedPingOutput)
		}

		// extract ms
		outms, err := strconv.Atoi(string(oout[1]))
		if err != nil {
			return -1, fmt.Errorf("error parsing time duration from cmd")
		}

		return time.Duration(outms), nil
	}
	NetDial LatencyFunc = func(url string) (time.Duration, error) {
		t1 := time.Now()
		_, err := net.Dial("tcp", url)
		if err != nil {
			return -1, fmt.Errorf("failed to dail address")
		}

		dur := time.Since(t1)
		return dur, nil
	}
)

type Latency struct {
	id      string
	servers []string
	latency map[string]time.Duration
	maxTime time.Duration
}

func NewLatencyChallenge(ctx context.Context) *Latency {
	return &Latency{id: "kjnfd", servers: Servers}
}

func RetrieveEndpoints(cts context.Context) ([]string, error) {
	return Servers, nil
}

func (l *Latency) Resolve(ctx context.Context, latencyFunc LatencyFunc) {
	l.latency = map[string]time.Duration{}

	for i := 0; i < len(l.servers); i++ {
		url := l.servers[i]
		dur, err := calcLatency(ctx, url, latencyFunc)
		if err != nil {
			log.Printf("failed while retriving latency information for url %s with error: %v\n", url, err)
			continue
		}
		l.latency[url] = dur
	}
}

func calcLatency(ctx context.Context, url string, latencyFunc LatencyFunc) (time.Duration, error) {

	dur, err := latencyFunc(url)
	return dur, err
}
