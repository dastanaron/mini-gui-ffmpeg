package main

import (
	"fmt"
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

	if err != nil {
		log.Fatal("Error:", err)
	}

	err = gui.AddFromFile(path.Join(rootPath, "./main.glade"))
	helpers.CheckError("Error", err)

	obj, err := gui.GetObject("window_main")
	helpers.CheckError("Error", err)

	win := obj.(*gtk.Window)
	win.SetTitle("GUI MINI FFMPEG")
	win.SetDefaultSize(500, 100)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	obj, _ = gui.GetObject("error_dialog")
	helpers.CheckError("Error", err)

	errorDialog := obj.(*gtk.MessageDialog)
	errorDialog.SetTitle("Error")
	errorDialog.SetDefaultSize(200, 100)

	obj, _ = gui.GetObject("choose_file_open")
	chooseFileOpen := obj.(*gtk.FileChooserButton)

	obj, _ = gui.GetObject("choose_file_save")
	chooseFileSave := obj.(*gtk.Button)
	obj, _ = gui.GetObject("button_run")
	buttonRun := obj.(*gtk.Button)

	obj, _ = gui.GetObject("file_save_dialog")
	fileSaveDialog := obj.(*gtk.FileChooserDialog)

	obj, _ = gui.GetObject("save_dialog_save")
	buttonSaveFile := obj.(*gtk.Button)

	obj, _ = gui.GetObject("save_dialog_cancel")
	buttonSaveFileCloseDialog := obj.(*gtk.Button)

	obj, _ = gui.GetObject("save_file_format")
	saveFileFormatCombox := obj.(*gtk.ComboBoxText)

	obj, _ = gui.GetObject("select_type")
	convertingType := obj.(*gtk.ComboBoxText)

	obj, _ = gui.GetObject("output_text")
	outputText := obj.(*gtk.TextView)

	var pathToSave string

	chooseFileSave.Connect("clicked", func() {
		fileSaveDialog.Show()
	})

	buttonSaveFileCloseDialog.Connect("clicked", func() {
		fileSaveDialog.Hide()
	})

	buttonSaveFile.Connect("clicked", func() {
		fileName := fileSaveDialog.FileChooser.GetFilename()
		format := saveFileFormatCombox.GetActiveText()

		if format == "*" {
			format = ""
		}

		pathToSave = fmt.Sprintf("%s%s", fileName, format)
		fileSaveDialog.Hide()
	})

	obj, _ = gui.GetObject("error_dialog_close")
	errorDialogCloseButton := obj.(*gtk.Button)
	errorDialogCloseButton.Connect("clicked", func() {
		errorDialog.Hide()
	})

	buttonRun.Connect("clicked", func() {
		filesToOpen, _ := chooseFileOpen.GetFilenames()

		outTextBuffer, err := outputText.GetBuffer()
		helpers.CheckError("Error receiving buffer", err)
		outTextBuffer.SetText(" ")

		if len(filesToOpen) == 0 {
			errorDialog.ShowAll()
			errorDialog.SetMarkup("Don't selected files")
			errorDialog.FormatSecondaryMarkup("%s", "You need to select file to save and to open")
			return
		}

		srcPath := filesToOpen[0]

		converter := ffmpeg.NewConverter(srcPath, pathToSave)

		convertingTypeId := convertingType.GetActiveID()

		switch convertingTypeId {
		case "0":
			converter.ConvertTelegram()
		case "1":
			converter.CutAudio()
		default:
			errorDialog.ShowAll()
			errorDialog.SetMarkup("Don't selected convert type")
			errorDialog.FormatSecondaryMarkup("%s", "You need to select convert type")
			return
		}

		report := fmt.Sprintf("%s", converter.CmdOutput)

		if report == "" {
			report = "Done without errors"
		}

		outTextBuffer.SetText(report)

	})

	win.ShowAll()

	gtk.Main()
}
