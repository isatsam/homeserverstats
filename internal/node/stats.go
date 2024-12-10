package node

import (
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type Stats struct {
	Name   string
	Time   string
	Cpu    float64
	Mem    float64
	Uptime string
}

func newStats(name string) Stats {
	var s Stats
	s.Name = name

	s.Time = time.Now().Format("15:04:05")

	cpuData, _ := cpu.Percent(0, false)
	s.Cpu = math.Round(cpuData[0])

	memory, _ := mem.VirtualMemory()
	s.Mem = math.Round(float64(memory.UsedPercent))

	uptime, _ := host.Uptime()
	convTime := time.Duration(uptime * uint64(time.Second)).Hours()
	days := int(convTime / 24)
	hours := int(convTime - (math.Round(convTime/24) * 24))

	var formattedUptime string
	if days > 0 {
		formattedUptime = fmt.Sprintf("%dd %dh", days, hours)
	} else {
		formattedUptime = fmt.Sprintf("%dh", hours)
	}
	s.Uptime = formattedUptime

	return s
}

func (s *Stats) print() {
	fmt.Printf("%s  CPU: %f  MEM: %f\n", s.Time, s.Cpu, s.Mem)
}
