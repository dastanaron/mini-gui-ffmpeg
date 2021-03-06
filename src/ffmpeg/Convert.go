package ffmpeg

import (
	"os/exec"
)

type Converter struct {
	ProgressChannel chan int
	Command         string
	InputFile       string
	OutputFile      string
	CmdOutput       []byte
}

func NewConverter(srcFilePath string, outputFilePath string) Converter {
	return Converter{
		ProgressChannel: make(chan int),
		Command:         "ffmpeg",
		InputFile:       srcFilePath,
		OutputFile:      outputFilePath,
	}
}

func (conv *Converter) ConvertTelegram() Converter {
	var args []string

	args = append(args, "-y", "-i", conv.InputFile)
	args = append(args, "-pix_fmt", "yuv420p")
	args = append(args, "-codec:a", "aac")
	args = append(args, "-c:v", "libx264")
	args = append(args, conv.OutputFile)

	cmd := exec.Command(conv.Command, args...)

	ProgressPipe(cmd, conv.ProgressChannel)

	return *conv
}

func (conv *Converter) CutAudio() Converter {
	var args []string

	args = append(args, "-y", "-i", conv.InputFile)
	args = append(args, "-vn", "-ar", "44100")
	args = append(args, "-ac", "2")
	args = append(args, "-ab", "192K")
	args = append(args, "-f", "mp3")
	args = append(args, conv.OutputFile)

	cmd := exec.Command(conv.Command, args...)

	ProgressPipe(cmd, conv.ProgressChannel)

	return *conv
}
