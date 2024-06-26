// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package helper

import (
	"ServeBin/data/response"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
)

// Gets the CPU load of the machine.
func getCPULoad() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatal(err)
	}
	return percent[0]
}

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

// Get the machine's Disk Usage
func getDiskUsage() (uint64, uint64, uint64) {
	usage, err := disk.Usage("/")
	if err != nil {
		log.Fatal(err)
	}
	return usage.Total, usage.Free, usage.Used
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
func getDiskOperationsAndPartitions() (int64, int64, uint64) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		log.Fatal(err)
	}

	var readCount int64
	var writeCount int64
	var partitionsCount uint64
	checkedDevices := make(map[string]bool)

	for _, partition := range partitions {
		usage, err := disk.IOCounters(partition.Device)
		if err != nil {
			log.Fatal(err)
		}

		if len(usage) > 0 {
			for device, value := range usage {
				if !checkedDevices[device] {
					readCount += int64(value.ReadCount)
					writeCount += int64(value.WriteCount)

					fmt.Println(usage)

					partitionsCount += 1
					checkedDevices[device] = true // Mark device as checked
				}
			}
		}
	}

	return readCount, writeCount, partitionsCount
}

// Get the server status on the basis of the different parameters
func GetHeartbeats() response.HeartbeatStats {
	stats := response.HeartbeatStats{}

	// Physical and Logical CPU Count
	stats.PhysicalAndLogicalCPUCount = runtime.NumCPU()

	// CPU Load
	stats.CPULoad = getCPULoad()

	// Memory Usage
	totalRAM, usedRAM := getRAMUsage()
	stats.RAM = response.RAMStats{
		TotalRAM: totalRAM,
		UsedRAM:  usedRAM,
	}

	// Disk Usage
	totalDisk, freeDisk, usedDisk := getDiskUsage()
	stats.Disk = response.DiskStats{
		TotalDiskSpace: totalDisk,
		FreeDiskSpace:  freeDisk,
		UsedDiskSpace:  usedDisk,
	}

	// Disk Usage
	readCount, writeCount, partitionsCount := getDiskOperationsAndPartitions()
	stats.Disk.ReadWrite = response.DiskReadWrite{
		Read:    readCount,
		Written: writeCount,
	}
	stats.Disk.Partitions = partitionsCount

	// Network Latency
	minLatency, avg, maxLatency := getNetworkLatency("https://ping.atishir.co", 5)
	stats.NetworkLatency = response.NetworkLatency{
		Min: minLatency,
		Avg: avg,
		Max: maxLatency,
	}

	return stats
}
