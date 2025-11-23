package handlers

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func init() {
	// Configure gopsutil to use host paths when running in Docker
	// These environment variables tell gopsutil where to find host system info
	if hostProc := os.Getenv("HOST_PROC"); hostProc != "" {
		os.Setenv("HOST_PROC", hostProc)
	}
	if hostSys := os.Getenv("HOST_SYS"); hostSys != "" {
		os.Setenv("HOST_SYS", hostSys)
	}
	if hostEtc := os.Getenv("HOST_ETC"); hostEtc != "" {
		os.Setenv("HOST_ETC", hostEtc)
	}
}

// SystemStats represents the system statistics
type SystemStats struct {
	Timestamp    string       `json:"timestamp"`
	Hostname     string       `json:"hostname"`
	Platform     string       `json:"platform"`
	OS           string       `json:"os"`
	Architecture string       `json:"architecture"`
	CPUCount     int          `json:"cpu_count"`
	CPU          CPUStats     `json:"cpu"`
	Memory       MemoryStats  `json:"memory"`
	Disk         []DiskStats  `json:"disk"`
	Network      NetworkStats `json:"network"`
	Uptime       UptimeStats  `json:"uptime"`
	Load         LoadStats    `json:"load"`
}

// CPUStats represents CPU statistics
type CPUStats struct {
	UsagePercent float64   `json:"usage_percent"`
	PerCoreUsage []float64 `json:"per_core_usage"`
	Temperature  float64   `json:"temperature,omitempty"`
	ModelName    string    `json:"model_name"`
}

// MemoryStats represents memory statistics
type MemoryStats struct {
	Total       uint64  `json:"total_bytes"`
	Available   uint64  `json:"available_bytes"`
	Used        uint64  `json:"used_bytes"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free_bytes"`
	TotalGB     float64 `json:"total_gb"`
	UsedGB      float64 `json:"used_gb"`
	AvailableGB float64 `json:"available_gb"`
}

// DiskStats represents disk statistics
type DiskStats struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mount_point"`
	FSType      string  `json:"fs_type"`
	Total       uint64  `json:"total_bytes"`
	Used        uint64  `json:"used_bytes"`
	Free        uint64  `json:"free_bytes"`
	UsedPercent float64 `json:"used_percent"`
	TotalGB     float64 `json:"total_gb"`
	UsedGB      float64 `json:"used_gb"`
	FreeGB      float64 `json:"free_gb"`
}

// NetworkStats represents network statistics
type NetworkStats struct {
	Interfaces []NetworkInterface `json:"interfaces"`
	TotalSent  uint64             `json:"total_sent_bytes"`
	TotalRecv  uint64             `json:"total_recv_bytes"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name      string   `json:"name"`
	Addresses []string `json:"addresses"`
	BytesSent uint64   `json:"bytes_sent"`
	BytesRecv uint64   `json:"bytes_recv"`
}

// UptimeStats represents system uptime
type UptimeStats struct {
	UptimeSeconds   uint64  `json:"uptime_seconds"`
	UptimeDays      float64 `json:"uptime_days"`
	UptimeFormatted string  `json:"uptime_formatted"`
	BootTime        string  `json:"boot_time"`
}

// LoadStats represents system load averages
type LoadStats struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// StatsHandler returns real-time system statistics
func StatsHandler(c *gin.Context) {
	stats := SystemStats{
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Architecture: runtime.GOARCH,
	}

	// Get hostname and platform info
	if hostInfo, err := host.Info(); err == nil {
		stats.Hostname = hostInfo.Hostname
		stats.Platform = hostInfo.Platform
		stats.OS = hostInfo.OS
		stats.Uptime = UptimeStats{
			UptimeSeconds:   hostInfo.Uptime,
			UptimeDays:      float64(hostInfo.Uptime) / 86400.0,
			UptimeFormatted: formatUptime(hostInfo.Uptime),
			BootTime:        time.Unix(int64(hostInfo.BootTime), 0).Format(time.RFC3339),
		}
	}

	// Get CPU info
	stats.CPUCount = runtime.NumCPU()
	if cpuPercent, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercent) > 0 {
		stats.CPU.UsagePercent = cpuPercent[0]
	}

	if perCorePercent, err := cpu.Percent(time.Second, true); err == nil {
		stats.CPU.PerCoreUsage = perCorePercent
	}

	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		stats.CPU.ModelName = cpuInfo[0].ModelName
	}

	// Get CPU temperature (may not work on all systems)
	if temps, err := host.SensorsTemperatures(); err == nil && len(temps) > 0 {
		// Try to find CPU temperature
		for _, temp := range temps {
			if temp.SensorKey == "coretemp_core_0" || temp.SensorKey == "cpu_thermal" {
				stats.CPU.Temperature = temp.Temperature
				break
			}
		}
	}

	// Get memory info
	if memInfo, err := mem.VirtualMemory(); err == nil {
		stats.Memory = MemoryStats{
			Total:       memInfo.Total,
			Available:   memInfo.Available,
			Used:        memInfo.Used,
			UsedPercent: memInfo.UsedPercent,
			Free:        memInfo.Free,
			TotalGB:     float64(memInfo.Total) / 1024 / 1024 / 1024,
			UsedGB:      float64(memInfo.Used) / 1024 / 1024 / 1024,
			AvailableGB: float64(memInfo.Available) / 1024 / 1024 / 1024,
		}
	}

	// Get disk info
	if partitions, err := disk.Partitions(false); err == nil {
		stats.Disk = make([]DiskStats, 0)
		for _, partition := range partitions {
			// Filter out virtual filesystems and bind mounts
			if isRealFilesystem(partition.Fstype, partition.Mountpoint) {
				if usage, err := disk.Usage(partition.Mountpoint); err == nil {
					diskStat := DiskStats{
						Device:      partition.Device,
						MountPoint:  partition.Mountpoint,
						FSType:      partition.Fstype,
						Total:       usage.Total,
						Used:        usage.Used,
						Free:        usage.Free,
						UsedPercent: usage.UsedPercent,
						TotalGB:     float64(usage.Total) / 1024 / 1024 / 1024,
						UsedGB:      float64(usage.Used) / 1024 / 1024 / 1024,
						FreeGB:      float64(usage.Free) / 1024 / 1024 / 1024,
					}
					stats.Disk = append(stats.Disk, diskStat)
				}
			}
		}
	}

	// Get network info
	stats.Network = NetworkStats{
		Interfaces: make([]NetworkInterface, 0),
	}

	if interfaces, err := net.Interfaces(); err == nil {
		for _, iface := range interfaces {
			addresses := make([]string, 0)
			for _, addr := range iface.Addrs {
				addresses = append(addresses, addr.Addr)
			}

			netInterface := NetworkInterface{
				Name:      iface.Name,
				Addresses: addresses,
			}

			stats.Network.Interfaces = append(stats.Network.Interfaces, netInterface)
		}
	}

	if ioCounters, err := net.IOCounters(false); err == nil && len(ioCounters) > 0 {
		stats.Network.TotalSent = ioCounters[0].BytesSent
		stats.Network.TotalRecv = ioCounters[0].BytesRecv
	}

	// Get load averages
	if loadAvg, err := load.Avg(); err == nil {
		stats.Load = LoadStats{
			Load1:  loadAvg.Load1,
			Load5:  loadAvg.Load5,
			Load15: loadAvg.Load15,
		}
	}

	c.JSON(http.StatusOK, stats)
}

// formatUptime formats uptime seconds into a human-readable string
func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, secs)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, secs)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, secs)
	}
	return fmt.Sprintf("%ds", secs)
}

// isRealFilesystem filters out virtual filesystems and bind mounts
func isRealFilesystem(fstype, mountpoint string) bool {
	// List of real filesystem types to include
	realFilesystems := map[string]bool{
		"ext4":  true,
		"ext3":  true,
		"ext2":  true,
		"xfs":   true,
		"btrfs": true,
		"ntfs":  true,
		"vfat":  true,
		"apfs":  true, // macOS
		"hfs":   true, // macOS
		"zfs":   true,
		"f2fs":  true,
	}

	// List of virtual/system filesystems to exclude
	excludeFilesystems := map[string]bool{
		"tmpfs":     true,
		"devtmpfs":  true,
		"devfs":     true,
		"sysfs":     true,
		"proc":      true,
		"devpts":    true,
		"cgroup":    true,
		"cgroup2":   true,
		"pstore":    true,
		"bpf":       true,
		"tracefs":   true,
		"debugfs":   true,
		"mqueue":    true,
		"hugetlbfs": true,
		"fusectl":   true,
		"fuse":      true,
		"overlay":   true,
		"squashfs":  true,
		"iso9660":   true,
	}

	// Exclude system mountpoints
	excludeMountpoints := map[string]bool{
		"/dev":           true,
		"/dev/shm":       true,
		"/run":           true,
		"/sys":           true,
		"/proc":          true,
		"/sys/fs/cgroup": true,
		"/boot/efi":      true,
	}

	// Check if filesystem type is explicitly excluded
	if excludeFilesystems[fstype] {
		return false
	}

	// Check if mountpoint is explicitly excluded
	if excludeMountpoints[mountpoint] {
		return false
	}

	// Check if it's a real filesystem type
	if realFilesystems[fstype] {
		return true
	}

	// Exclude anything that looks like a file (not a directory mount)
	// This catches Docker bind mounts like /etc/resolv.conf
	if len(mountpoint) > 0 && mountpoint[0] == '/' {
		// If it contains a file extension or looks like a file path
		for i := len(mountpoint) - 1; i >= 0; i-- {
			if mountpoint[i] == '/' {
				break
			}
			if mountpoint[i] == '.' {
				return false // Likely a file, not a directory
			}
		}
	}

	return false
}
