package main

import (
	"fmt"
	controlgui "gui-feed-analyzer/control-gui"
	"gui-feed-analyzer/ffmpeg"
	"gui-feed-analyzer/helpers"
	"log"
	"os"
	"path"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	gui, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Error:", err)
	}

	rootPath, err := os.Getwd()

	helpers.CheckError("Error", err)

	err = gui.AddFromFile(path.Join(rootPath, "./main.glade"))
	helpers.CheckError("Error", err)

	controlgui, err := controlgui.NewGUIController(gui)
	helpers.CheckError("Error", err)

	controlgui.Common.MainWindow.SetTitle("GUI MINI FFMPEG")
	controlgui.Common.MainWindow.SetDefaultSize(500, 100)
	controlgui.Common.MainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})

	var pathToSave string

	controlgui.SaveDialog.OpenFileSaveButton.Connect("clicked", func() {
		controlgui.SaveDialog.FileSaveDialog.Show()
	})

	controlgui.SaveDialog.CancelButton.Connect("clicked", func() {
		controlgui.SaveDialog.FileSaveDialog.Hide()
	})

	controlgui.SaveDialog.SaveButton.Connect("clicked", func() {
		fileName := controlgui.SaveDialog.FileSaveDialog.FileChooser.GetFilename()
		format := controlgui.SaveDialog.FileFormatBox.GetActiveText()

		if format == "*" {
			format = ""
		}

		pathToSave = fmt.Sprintf("%s%s", fileName, format)
		controlgui.SaveDialog.FileSaveDialog.Hide()
	})

	controlgui.ErrorDialog.CloseButton.Connect("clicked", func() {
		controlgui.ErrorDialog.ErrorDialog.Hide()
	})

	controlgui.Common.StartButton.Connect("clicked", func() {
		filesToOpen, _ := controlgui.Common.OpenFileButton.GetFilenames()

		outTextBuffer, err := controlgui.Common.ReportOutput.GetBuffer()
		helpers.CheckError("Error receiving buffer", err)
		outTextBuffer.SetText(" ")

		if len(filesToOpen) == 0 {
			controlgui.ErrorDialog.ErrorDialog.ShowAll()
			controlgui.ErrorDialog.ErrorDialog.SetMarkup("Don't selected files")
			controlgui.ErrorDialog.ErrorDialog.FormatSecondaryMarkup("%s", "You need to select file to save and to open")
			return
		}

		srcPath := filesToOpen[0]

		converter := ffmpeg.NewConverter(srcPath, pathToSave)

		convertingTypeId := controlgui.Common.ConvertingTypeBox.GetActiveID()

		switch convertingTypeId {
		case "0":
			converter.ConvertTelegram()
		case "1":
			converter.CutAudio()
		default:
			controlgui.ErrorDialog.ErrorDialog.ShowAll()
			controlgui.ErrorDialog.ErrorDialog.SetMarkup("Don't selected convert type")
			controlgui.ErrorDialog.ErrorDialog.FormatSecondaryMarkup("%s", "You need to select convert type")
			return
		}

		report := fmt.Sprintf("%s", converter.CmdOutput)

		if report == "" {
			report = "Done without errors"
		}

		outTextBuffer.SetText(report)

	})

	controlgui.Common.MainWindow.ShowAll()

	gtk.Main()
}
