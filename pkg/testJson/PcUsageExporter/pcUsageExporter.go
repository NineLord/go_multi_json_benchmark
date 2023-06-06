package PcUsageExporter

import (
	"fmt"
	"github.com/struCoder/pidusage"
	"os"
	"time"
)

type PcUsage struct {
	Cpu float64
	Ram float64
}

func Main(senderToMain chan []PcUsage, receiveFromMain chan bool, sampleInterval uint) {
	sampleIntervalDuration := time.Duration(sampleInterval)
	result := make([]PcUsage, 0)
	for {
		select {
		case <-receiveFromMain:
			senderToMain <- result
			close(senderToMain)
			return
		default:
			systemInfo, err := pidusage.GetStat(os.Getpid())
			if err != nil {
				panic(fmt.Sprintf("Failed to get PC Usage information: %s", err))
			}
			result = append(result, PcUsage{
				Cpu: systemInfo.CPU,
				Ram: (systemInfo.Memory / 1024) / 1024,
			})
			time.Sleep(sampleIntervalDuration * time.Millisecond)
		}
	}
}
