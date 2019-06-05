package utils

import (
	"bytes"
	"os/exec"
)

//执行系统命令，返回执行结果
//同步操作
func Exec(cmd string, shell bool) (string, error) {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			return "", err
		}
		return string(out), err
	} else {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			return "", err
		}
		return string(out), err
	}
}

const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}
