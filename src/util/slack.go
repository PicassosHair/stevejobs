package util

import (
	"fmt"
	"os/exec"
)

// LogSlack send a message to slack channel.
func LogSlack(messageType string, message string) {
  _, err := exec.Command("/usr/src/app/jobs/log_slack.sh", messageType, message).Output()

  if err != nil {
    fmt.Println(err.Error())
  }
}