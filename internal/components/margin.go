package components

import "fyne.io/fyne/v2"

type marginLayout struct {
	margin fyne.Size
}

func NewMarginLayout(margin fyne.Size) fyne.Layout {
	return &marginLayout{margin: margin}
}

func (m *marginLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, obj := range objects {
		obj.Resize(size.Subtract(m.margin).Max(fyne.NewSize(0, 0)))
		obj.Move(fyne.NewPos(m.margin.Width/2, m.margin.Height/2))
	}
}

func (m *marginLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, obj := range objects {
		minSize = minSize.Max(obj.MinSize())
	}
	return minSize.Add(m.margin)
}
