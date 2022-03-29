package mpv

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/logrusorgru/aurora/v3"
)

func Play(video string, title string) {
	cmd := exec.Command("mpv", video, "--force-media-title="+title)
	if err := cmd.Start(); err != nil {
		fmt.Println(aurora.Red(err))
		os.Exit(1)
	}
}
