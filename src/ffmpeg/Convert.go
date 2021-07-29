package ffmpeg

import (
	"fmt"
	"os/exec"
)

func ConvertTelegram(srcFilePath string, outputFilePath string) []byte {
	var args []string

	command := "ffmpeg"

	args = append(args, "-y", "-i", srcFilePath)
	args = append(args, "-pix_fmt", "yuv420p")
	args = append(args, "-codec:a", "aac")
	args = append(args, "-c:v", "libx264")
	args = append(args, outputFilePath)

	cmd := exec.Command(command, args...)

	output, _ := cmd.CombinedOutput()

	fmt.Printf("%s", output)

	return output
}
