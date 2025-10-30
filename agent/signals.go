package agent

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (sa *SysAgent) MsgUsers(msgStr string) error {
	fmt.Printf("Message %v users\n", len(sa.users))
	for userID, _ := range sa.users {
		fmt.Printf("Message user %v\n", userID)
		msgTel := tgbotapi.NewMessage(userID, msgStr)
		if _, err := sa.bot.Send(msgTel); err != nil {
			return fmt.Errorf("Failed to send message %v: error %v\n", err)
		}
	}
	return nil
}

type GPUStats struct {
	MemUsed  int
	MemTotal int
	Temp     int
}

func parseGPUInfo(text string) (GPUStats, error) {
	g := GPUStats{}
	// extract temperature
	tempRe := regexp.MustCompile(`\|\s*\d+%[\s]+(\d+)C`)
	if match := tempRe.FindStringSubmatch(text); len(match) > 1 {
		g.Temp, _ = strconv.Atoi(match[1])
	}
	// extract video memory: “xxxxxMiB / yyyyyMiB”
	memRe := regexp.MustCompile(`\|\s*(\d+)MiB\s*/\s*(\d+)MiB\s*\|`)
	if match := memRe.FindStringSubmatch(text); len(match) > 2 {
		g.MemUsed, _ = strconv.Atoi(match[1])
		g.MemTotal, _ = strconv.Atoi(match[2])
	}

	return g, nil
}

func (sa *SysAgent) CheckNvidiaGpus() {
	type GpuItem struct {
		Name    string
		MemMin  int
		MemMax  int
		TempMin int
		TempMax int
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error fetching hostname:", err)
	}
	// initial message when bot starts
	if err := sa.MsgUsers(
		fmt.Sprintf("\U000026A0 Bot started for %v at %v",
			hostname, time.Now().Format("2006-01-02 15:04:05"))); err != nil {
		fmt.Printf("Error send message %v", err)
	}

	for {
		for gpuIdx, gpuItem := range []GpuItem{
			{
				// id:0
				Name:    "RTX 5090",
				MemMin:  2000,
				MemMax:  28000,
				TempMin: 30,
				TempMax: 65,
			}, {
				// id:1
				Name:    "RTX 5060 Ti",
				MemMin:  2000,
				MemMax:  12000,
				TempMin: 35,
				TempMax: 75,
			},
		} {
			output, err := sysExecutor("nvidia-smi", "-i",
				fmt.Sprintf("%d", gpuIdx))
			if err != nil {
				time.Sleep(time.Minute)
				continue
			}
			gpuStats, err := parseGPUInfo(output)
			if err != nil {
				err := sa.MsgUsers(
					fmt.Sprintf("\U000026A0 Error parse GPU info %v", err))
				if err != nil {
					fmt.Printf("Error send message %v", err)
				}
				time.Sleep(time.Minute)
				continue
			}

			if gpuStats.MemUsed < gpuItem.MemMin {
				err := sa.MsgUsers(
					fmt.Sprintf("\U000026A0 Memory usage too low for %v: %vMB",
						gpuItem.Name, gpuStats.MemUsed))
				if err != nil {
					fmt.Printf("Error send message %v", err)
				}
			} else if gpuStats.MemUsed > gpuItem.MemMax {
				err := sa.MsgUsers(
					fmt.Sprintf("\U000026A0 Memory usage too high for %v: %vMB",
						gpuItem.Name, gpuStats.MemUsed))
				if err != nil {
					fmt.Printf("Error send message %v", err)
				}
			}

			if gpuStats.Temp < gpuItem.TempMin {
				err := sa.MsgUsers(
					fmt.Sprintf("\U000026A0 Temperature too low for %v: %vC",
						gpuItem.Name, gpuStats.Temp))
				if err != nil {
					fmt.Printf("Error send message %v", err)
				}
			} else if gpuStats.Temp > gpuItem.TempMax {
				err := sa.MsgUsers(
					fmt.Sprintf("\U000026A0 Temperature too high for %v: %vC",
						gpuItem.Name, gpuStats.Temp))
				if err != nil {
					fmt.Printf("Error send message %v", err)
				}
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
