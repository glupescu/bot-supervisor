package agent

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"bot-supervisor/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ExecData struct {
	trim    int
	cmd     string
	args    []string
	sysFunc func(cmd string, arg ...string) (string, error)
}

type SysAgent struct {
	name         string
	bot          *tgbotapi.BotAPI
	execFuncDict map[string]ExecData
	execWords    []string
	sTime        time.Time
}

func (sa *SysAgent) Init(name string, bot *tgbotapi.BotAPI) error {
	sa.name = name
	sa.bot = bot
	sa.sTime = time.Now()
	sa.execFuncDict = map[string]ExecData{
		"temp": {
			cmd:     "sensors",
			args:    []string{},
			sysFunc: sysExecutor,
		},
		"top": {
			trim:    256,
			cmd:     "top",
			args:    []string{"-b", "-n", "1"},
			sysFunc: sysExecutor,
		},
		"cpu": {
			trim:    256,
			cmd:     "top",
			args:    []string{"-b", "-n", "1"},
			sysFunc: sysExecutor,
		},
		"mem": {
			cmd:     "free",
			sysFunc: sysExecutor,
		},
		"speedtest": {
			cmd:     "speedtest-cli",
			sysFunc: sysExecutor,
		},
		"disk": {
			cmd:     "df",
			args:    []string{"-h"},
			sysFunc: sysExecutor,
		},
		"gpu": {
			cmd:     "nvidia-smi",
			args:    []string{},
			sysFunc: sysExecutor,
		},
	}
	sa.execFuncDict["help"] = ExecData{
		cmd:  "help",
		args: []string{},
		sysFunc: func(_ string, _ ...string) (string, error) {
			helpStr := "Available commands:\n"
			for execKey, _ := range sa.execFuncDict {
				sa.execWords = append(sa.execWords, execKey)
				helpStr += fmt.Sprintf("\t%v\n", execKey)
			}
			return helpStr, nil
		},
	}
	return nil
}

func (sa *SysAgent) GetBot() *tgbotapi.BotAPI {
	return sa.bot
}

func (sa *SysAgent) Serve(msg string,
	userRole user.Type,
	chatID int64) (string, error) {
	if userRole != user.FullAccess {
		return "", fmt.Errorf("user role is not allowed")
	}
	msg = strings.ToLower(msg)
	for key, ed := range sa.execFuncDict {
		if strings.Contains(msg, key) {
			_, err := sa.bot.Send(tgbotapi.NewMessage(chatID,
				fmt.Sprintf("Will execute %v %v...", ed.cmd, ed.args)))
			if err != nil {
				return "Error while sending message", err
			}
			rsp, err := ed.sysFunc(ed.cmd, ed.args...)
			if err != nil {
				return "", err
			}
			if ed.trim > 0 {
				trim := min(ed.trim, len(rsp))
				return rsp[:trim], nil
			}
			return rsp, nil
		}
	}

	return "", nil
}
func (sa *SysAgent) Characteristics() []string {
	return sa.execWords
}

func sysExecutor(cmd string, arg ...string) (string, error) {
	exeCmd := exec.Command(cmd, arg...)

	// Run the command and capture the output
	output, err := exeCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run cmd: %v", err)
	}

	// Print the output as a string
	return string(output), nil
}

func (sa *SysAgent) IsSysLibrary(msg string, userRole user.Type) bool {
	if userRole != user.FullAccess {
		return false
	}
	for k, _ := range sa.execFuncDict {
		if strings.Contains(msg, k) {
			return true
		}
	}
	return false
}

// TrySysLibrary attempts a sys resolve
func (sa *SysAgent) TrySysLibrary(msg string, userRole user.Type) (string, error) {
	if userRole != user.FullAccess {
		return "", fmt.Errorf("user role is not allowed")
	}
	for key, ed := range sa.execFuncDict {
		if strings.Contains(msg, key) {
			rsp, err := ed.sysFunc(ed.cmd, ed.args...)
			if err != nil {
				return "", err
			}
			if ed.trim > 0 {
				trim := min(ed.trim, len(rsp))
				return rsp[:trim], nil
			}
			return rsp, nil
		}
	}
	return "", errors.New("system library not found")
}
