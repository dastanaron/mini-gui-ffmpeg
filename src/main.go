package main

import (
	_ "embed"
	"errors"
	"fmt"
	"gui-mini-ffmpeg/ffmpeg"
	guiController "gui-mini-ffmpeg/gui-controller"
	"gui-mini-ffmpeg/helpers"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

//go:embed main.glade
var gladeInterface string

//go:embed icon.png
var icon []byte

type ConvertingObject struct {
	InputFile      string
	OutputFile     string
	ConvertingType string
}

func main() {
	gtk.Init(nil)

	pixBuf, err := gdk.PixbufNewFromBytesOnly(icon)
	helpers.CheckError("Error", err)

	gui, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Error:", err)
	}

	err = gui.AddFromString(gladeInterface)
	helpers.CheckError("Error", err)

	controllerGUI, err := guiController.NewGUIController(gui)
	helpers.CheckError("Error", err)

	controllerGUI.Common.MainWindow.SetTitle("GUI MINI FFMPEG")
	controllerGUI.Common.MainWindow.SetDefaultSize(500, 100)
	controllerGUI.Common.MainWindow.SetIcon(pixBuf)
	controllerGUI.Common.MainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})

	convertingObject := &ConvertingObject{}

	handleSaveDialog(controllerGUI, convertingObject)

	controllerGUI.ErrorDialog.CloseButton.Connect("clicked", func() {
		controllerGUI.ErrorDialog.ErrorDialog.Hide()
	})

	controllerGUI.Common.MainWindow.ShowAll()

	handleRunButton(controllerGUI, convertingObject)

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
			helpers.CheckGUIError(controllerGUI, "Don't selected files", errors.New("you need to select file to save and to open"))
			return
		}

		convertingObject.InputFile = filesToOpen[0]

		converter := ffmpeg.NewConverter(convertingObject.InputFile, convertingObject.OutputFile)

		controllerGUI.Common.StartButton.SetSensitive(false)

		go func() {
			controllerGUI.Common.ProgressBar.SetFraction(0.00)
			for progressPercent := range converter.ProgressChannel {
				fraction := float64(progressPercent) / 100
				fmt.Println(float64(fraction))
				controllerGUI.Common.ProgressBar.SetFraction(fraction)
			}
			outTextBuffer.SetText("Done")
			controllerGUI.Common.StartButton.SetSensitive(true)
		}()

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
	})
}
