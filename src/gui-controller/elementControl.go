package GUIController

import (
	"github.com/gotk3/gotk3/gtk"
)

type GUICommonWindow struct {
	MainWindow        *gtk.Window
	OpenFileButton    *gtk.FileChooserButton
	StartButton       *gtk.Button
	ConvertingTypeBox *gtk.ComboBoxText
	ReportOutput      *gtk.TextView
}

type GUIErrorDialog struct {
	ErrorDialog *gtk.MessageDialog
	CloseButton *gtk.Button
}

type GUISaveFileDiaog struct {
	OpenFileSaveButton *gtk.Button
	FileSaveDialog     *gtk.FileChooserDialog
	SaveButton         *gtk.Button
	CancelButton       *gtk.Button
	FileFormatBox      *gtk.ComboBoxText
}

type GUIInterface struct {
	Common      GUICommonWindow
	SaveDialog  GUISaveFileDiaog
	ErrorDialog GUIErrorDialog
}

func NewGUIController(gtkBuilder *gtk.Builder) (*GUIInterface, error) {

	commonWindow := GUICommonWindow{}
	saveDialog := GUISaveFileDiaog{}
	errorDialog := GUIErrorDialog{}

	obj, err := gtkBuilder.GetObject("window_main")
	if err != nil {
		return nil, err
	}
	commonWindow.MainWindow = obj.(*gtk.Window)

	obj, err = gtkBuilder.GetObject("choose_file_open")
	if err != nil {
		return nil, err
	}
	commonWindow.OpenFileButton = obj.(*gtk.FileChooserButton)

	obj, err = gtkBuilder.GetObject("button_run")
	if err != nil {
		return nil, err
	}
	commonWindow.StartButton = obj.(*gtk.Button)

	obj, err = gtkBuilder.GetObject("select_type")
	if err != nil {
		return nil, err
	}
	commonWindow.ConvertingTypeBox = obj.(*gtk.ComboBoxText)

	obj, err = gtkBuilder.GetObject("output_text")
	if err != nil {
		return nil, err
	}
	commonWindow.ReportOutput = obj.(*gtk.TextView)

	obj, err = gtkBuilder.GetObject("choose_file_save")
	if err != nil {
		return nil, err
	}
	saveDialog.OpenFileSaveButton = obj.(*gtk.Button)

	obj, err = gtkBuilder.GetObject("file_save_dialog")
	if err != nil {
		return nil, err
	}
	saveDialog.FileSaveDialog = obj.(*gtk.FileChooserDialog)

	obj, err = gtkBuilder.GetObject("save_dialog_save")
	if err != nil {
		return nil, err
	}
	saveDialog.SaveButton = obj.(*gtk.Button)

	obj, err = gtkBuilder.GetObject("save_dialog_cancel")
	if err != nil {
		return nil, err
	}
	saveDialog.CancelButton = obj.(*gtk.Button)

	obj, err = gtkBuilder.GetObject("save_file_format")
	if err != nil {
		return nil, err
	}
	saveDialog.FileFormatBox = obj.(*gtk.ComboBoxText)

	obj, err = gtkBuilder.GetObject("error_dialog")
	if err != nil {
		return nil, err
	}
	errorDialog.ErrorDialog = obj.(*gtk.MessageDialog)

	obj, err = gtkBuilder.GetObject("error_dialog_close")
	if err != nil {
		return nil, err
	}
	errorDialog.CloseButton = obj.(*gtk.Button)

	guiInterface := &GUIInterface{
		Common:      commonWindow,
		SaveDialog:  saveDialog,
		ErrorDialog: errorDialog,
	}

	return guiInterface, nil
}
