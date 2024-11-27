package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("URL Shortener")
	w.Resize(fyne.NewSize(400, 200))

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter URL here")

	output := widget.NewLabel("Shortened URL will appear here")

	button := widget.NewButton("Shorten URL", func() {
		url := input.Text
		data, _ := json.Marshal(map[string]string{"url": url})

		resp, err := http.Post("http://localhost:8080/shorten", "application/json", bytes.NewBuffer(data))
		if err != nil {
			output.SetText("Error connecting to server")
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var result map[string]string
		json.Unmarshal(body, &result)

		if short, ok := result["short_url"]; ok {
			output.SetText(short)
		} else {
			output.SetText("Error shortening URL")
		}
	})

	content := container.NewVBox(
		widget.NewLabel("URL Shortener"),
		input,
		button,
		output,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
