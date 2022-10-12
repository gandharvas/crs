package main

import (
	"image/color"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/gandharvas/crs/internal"
)

const KuteGoAPIURL = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRcaeVZK-_yozg3QQwnXStPkcIUvvAxbf-vUw&usqp=CAU"

var crs int
var crsFileURL string

func acceptUserScore(estimateITA *widget.Label) *fyne.Container {

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter your CRS score")

	content := container.NewVBox(input, widget.NewButton("Calculate", func() {
		log.Println("Content was:", input.Text)
		var err error
		crs, err = strconv.Atoi(input.Text)
		if err != nil {
			log.Println("Invalid CRS")
		}
		predObj := new(internal.CRS)
		log.Println("File URL is:", crsFileURL)
		predObj.Get_crs_distribution(crsFileURL)

		_, ita := internal.Predict(predObj, crs)

		estimateITA.Text = "Your estimated ITA date is: " + ita.Format("DD-MM-YYYY")
		estimateITA.Refresh()
	}))

	return content
}

func listPreviousDates() fyne.Widget {
	// Get the latest dates
	internal.DownloadCRSDates()
	datesMap := internal.GetCRSDates()

	dates := make([]string, len(datesMap))
	i := 0
	for date, _ := range datesMap {
		dates[i] = date
		i++
	}
	//list := widget.NewList(
	//	func() int {
	//		return len(dates)
	//	},
	//	func() fyne.CanvasObject {
	//		return widget.NewLabel("template")
	//	},
	//	func(i widget.ListItemID, o fyne.CanvasObject) {
	//		o.(*widget.Label).SetText(dates[i])
	//		log.Println(dates[i])
	//		crsFileURL = datesMap[dates[i]]
	//	})

	combo := widget.NewSelect(dates, func(value string) {
		log.Println("Select set to", value)
		crsFileURL = datesMap[value]
	})

	return combo

}

func prepareLayout(app fyne.App, window fyne.Window) {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Quit", func() { app.Quit() }),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			dialog.ShowCustom("About", "Close", container.NewVBox(
				widget.NewLabel("Welcome to CRS Predictor. A simple app which calculates an estimated date for getting ITA"),
				widget.NewLabel("Version: v1-alpha1"),
				widget.NewLabel("Author: Gandharva Shankara Murthy"),
			), window)
		}))
	mainMenu := fyne.NewMainMenu(
		fileMenu,
		helpMenu,
	)
	window.SetMainMenu(mainMenu)
}

func highLevelText() fyne.CanvasObject {
	//Define a welcome text centered
	text := canvas.NewText("CRS score Predictor", color.Black)
	text.Alignment = fyne.TextAlignCenter
	return text
}

func displayImage() fyne.CanvasObject {
	var resource, _ = fyne.LoadResourceFromURLString(KuteGoAPIURL)
	canadaImage := canvas.NewImageFromResource(resource)
	canadaImage.SetMinSize(fyne.Size{Width: 500, Height: 500}) // by default size is 0, 0
	return canadaImage
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("CRS Predictor by Gandharva")

	prepareLayout(myApp, myWindow)
	text := highLevelText()

	canadaImage := displayImage()

	dateSeletionText := canvas.NewText("Please select a date of previous draw for prediction", color.Black)
	dateSeletionText.Alignment = fyne.TextAlignCenter
	dates := listPreviousDates()

	itaDateText := canvas.NewText("Your expected ITA date is", color.Black)
	itaDateText.Alignment = fyne.TextAlignCenter

	estimateITA := widget.NewLabel("")
	userScore := acceptUserScore(estimateITA)

	// Display a vertical box containing text, image and button
	box := container.NewVBox(
		text,
		canadaImage,
		dateSeletionText,
		dates,
		userScore,
		estimateITA,
	)

	// Display our content
	myWindow.SetContent(box)

	// Close the App when Escape key is pressed
	myWindow.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {

		if keyEvent.Name == fyne.KeyEscape {
			myApp.Quit()
		}
	})

	// Show window and run app
	myWindow.ShowAndRun()
}
