package CustomUITools

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
)

func NewWSpacer(w float32) *canvas.Rectangle {
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(w, 0))

	return spacer
}

func NewHSpacer(h float32) *canvas.Rectangle {
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(0, h))

	return spacer
}

func NewHWSpacer(h float32, w float32) *canvas.Rectangle {
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(w, h))

	return spacer
}

func NewColorHSpacer(h float32, color color.Color) *canvas.Rectangle {
	spacer := canvas.NewRectangle(color)
	spacer.SetMinSize(fyne.NewSize(0, h))

	return spacer
}

func NewColorWSpacer(w float32, color color.Color) *canvas.Rectangle {
	spacer := canvas.NewRectangle(color)
	spacer.SetMinSize(fyne.NewSize(w, 0))

	return spacer
}

func NewColorHWSpacer(h float32, w float32, color color.Color) *canvas.Rectangle {
	spacer := canvas.NewRectangle(color)
	spacer.SetMinSize(fyne.NewSize(w, h))

	return spacer
}

func NewSeparator() *canvas.Line {
	line := canvas.NewLine(color.RGBA{R: 44, G: 44, B: 46, A: 255})
	line.StrokeWidth = 2.5

	return line
}
