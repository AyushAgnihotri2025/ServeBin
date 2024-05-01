// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package helper

import (
	"ServeBin/data/response"
	"github.com/shirou/gopsutil/disk"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

// Gets and returns the Network Latency of the Host
func getNetworkLatency(url string, numAttempts int) (float64, float64, float64) {
	var latencies []float64

	for i := 0; i < numAttempts; i++ {
		startTime := time.Now()
		_, err := http.Get(url)
		endTime := time.Now()

		if err == nil {
			latency := endTime.Sub(startTime).Seconds() * 1000
			latencies = append(latencies, latency)
		} else {
			return 0, 0, 0
		}
	}

	if len(latencies) == 0 {
		return 0, 0, 0
	}

	minLatency := latencies[0]
	maxLatency := latencies[0]
	var sumLatency float64

	for _, latency := range latencies {
		if latency < minLatency {
			minLatency = latency
		}
		if latency > maxLatency {
			maxLatency = latency
		}
		sumLatency += latency
	}

	avgLatency := sumLatency / float64(len(latencies))

	return minLatency, avgLatency, maxLatency
}

// Get the OS current RAM usage
func getRAMUsage() (float64, float64) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	totalRAM := float64(memStats.TotalAlloc) / (1024 * 1024 * 1024) // Convert to gigabytes (GB)
	usedRAM := float64(memStats.Alloc) / (1024 * 1024 * 1024)       // Convert to gigabytes (GB)

	return totalRAM, usedRAM
}

// Get the OS disk related stuffs
func getDiskOperations() (uint64, uint64) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		log.Fatal(err)
	}

	var readCount uint64
	var writeCount uint64

	for _, partition := range partitions {
		usage, err := disk.IOCounters(partition.Mountpoint)
		if err != nil {
			log.Fatal(err)
		}

		readCount += usage[strconv.Itoa(0)].ReadCount
		writeCount += usage[strconv.Itoa(0)].WriteCount
	}

	return readCount, writeCount
}

// Get the server status on the basis of the different parameters
func GetHeartbeats() response.HeartbeatStats {
	stats := response.HeartbeatStats{}

	// Physical and Logical CPU Count
	stats.PhysicalAndLogicalCPUCount = runtime.NumCPU()

	// Memory Usage
	totalRAM, usedRAM := getRAMUsage()
	stats.RAM = response.RAMStats{
		TotalRAM: totalRAM,
		UsedRAM:  usedRAM,
	}

	// Disk Usage
	readCount, writeCount := getDiskOperations()
	stats.Disk = response.DiskStats{
		ReadWrite: response.DiskReadWrite{
			Read:    readCount,
			Written: writeCount,
		},
	}

	// Network Latency
	minLatency, avg, maxLatency := getNetworkLatency("https://ping.atishir.co", 5)
	stats.NetworkLatency = response.NetworkLatency{
		Min: minLatency,
		Avg: avg,
		Max: maxLatency,
	}

	return stats
}