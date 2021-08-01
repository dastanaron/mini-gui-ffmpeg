package main

import (
	"errors"
	"fmt"
	"gui-feed-analyzer/ffmpeg"
	guiController "gui-feed-analyzer/gui-controller"
	"gui-feed-analyzer/helpers"
	"log"
	"os"
	"path"

	"github.com/gotk3/gotk3/gtk"
)

type ConvertingObject struct {
	InputFile      string
	OutputFile     string
	ConvertingType string
}

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

	controllerGUI, err := guiController.NewGUIController(gui)
	helpers.CheckError("Error", err)

	controllerGUI.Common.MainWindow.SetTitle("GUI MINI FFMPEG")
	controllerGUI.Common.MainWindow.SetDefaultSize(500, 100)
	controllerGUI.Common.MainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})

	convertingObject := &ConvertingObject{}

	handleSaveDialog(controllerGUI, convertingObject)

	controllerGUI.ErrorDialog.CloseButton.Connect("clicked", func() {
		controllerGUI.ErrorDialog.ErrorDialog.Hide()
	})

	handleRunButton(controllerGUI, convertingObject)

	controllerGUI.Common.MainWindow.ShowAll()

	gtk.Main()
}

func handleSaveDialog(controllerGUI *guiController.GUIInterface, convertingObject *ConvertingObject) {
	controllerGUI.SaveDialog.OpenFileSaveButton.Connect("clicked", func() {
		controllerGUI.SaveDialog.FileSaveDialog.Show()
	})

	controllerGUI.SaveDialog.CancelButton.Connect("clicked", func() {
		controllerGUI.SaveDialog.FileSaveDialog.Hide()
	})

	controllerGUI.SaveDialog.SaveButton.Connect("clicked", func() {
		fileName := controllerGUI.SaveDialog.FileSaveDialog.FileChooser.GetFilename()
		format := controllerGUI.SaveDialog.FileFormatBox.GetActiveText()

		if format == "*" {
			format = ""
		}

		convertingObject.OutputFile = fmt.Sprintf("%s%s", fileName, format)
		controllerGUI.SaveDialog.FileSaveDialog.Hide()
	})
}

func handleRunButton(controllerGUI *guiController.GUIInterface, convertingObject *ConvertingObject) {
	controllerGUI.Common.StartButton.Connect("clicked", func() {
		filesToOpen, err := controllerGUI.Common.OpenFileButton.GetFilenames()

		if err != nil {
			helpers.CheckGUIError(controllerGUI, "Cannot open file", err)
			return
		}

		outTextBuffer, err := controllerGUI.Common.ReportOutput.GetBuffer()
		if err != nil {
			helpers.CheckGUIError(controllerGUI, "Error receiving buffer", err)
			return
		}

		outTextBuffer.SetText(" ")

		if len(filesToOpen) == 0 {
			helpers.CheckGUIError(controllerGUI, "Don't selected files", errors.New("You need to select file to save and to open"))
			return
		}

		convertingObject.InputFile = filesToOpen[0]

		converter := ffmpeg.NewConverter(convertingObject.InputFile, convertingObject.OutputFile)

		convertingObject.ConvertingType = controllerGUI.Common.ConvertingTypeBox.GetActiveID()

		switch convertingObject.ConvertingType {
		case "0":
			converter.ConvertTelegram()
		case "1":
			converter.CutAudio()
		default:
			controllerGUI.ErrorDialog.ErrorDialog.ShowAll()
			controllerGUI.ErrorDialog.ErrorDialog.SetMarkup("Don't selected convert type")
			controllerGUI.ErrorDialog.ErrorDialog.FormatSecondaryMarkup("%s", "You need to select convert type")
			return
		}

		report := fmt.Sprintf("%s", converter.CmdOutput)

		if report == "" {
			report = "Done without errors"
		}

		outTextBuffer.SetText(report)

	})
}
