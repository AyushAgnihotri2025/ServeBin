// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package response

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type AboutResponse struct {
	Version    string `json:"version"`
	ServerTime string `json:"serverTime"`
	Developer  string `json:"developer"`
	Contact    string `json:"contact"`
	SourceCode string `json:"sourceCode"`
}

type HeartbeatResponse struct {
	Stats  HeartbeatStats `json:"stats"`
	Status string         `json:"status"`
}

type HeartbeatStats struct {
	CPULoad                    float64        `json:"cpu_load"`
	Disk                       DiskStats      `json:"disk"`
	NetworkLatency             NetworkLatency `json:"network_latency"`
	PhysicalAndLogicalCPUCount int            `json:"physical_and_logical_cpu_count"`
	RAM                        RAMStats       `json:"ram"`
}

type DiskStats struct {
	FreeDiskSpace  uint64        `json:"free_disk_space"`
	ReadWrite      DiskReadWrite `json:"read_write"`
	TotalDiskSpace uint64        `json:"total_disk_space"`
	UsedDiskSpace  uint64        `json:"used_disk_space"`
	Partitions     uint64        `json:"partitions"`
}

type DiskReadWrite struct {
	Read    int64 `json:"read"`
	Written int64 `json:"written"`
}

type NetworkLatency struct {
	Min float64 `json:"min"`
	Avg float64 `json:"avg"`
	Max float64 `json:"max"`
}

type RAMStats struct {
	TotalRAM float64 `json:"total_ram"`
	UsedRAM  float64 `json:"used_ram"`
}
