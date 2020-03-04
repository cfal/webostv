package main

type volume struct {
	*Slider
}

func newVolume() *volume {
	w := NewSlider()
	w.SetBorder(true)
	w.SetTitle("Volume")
	return &volume{w}
}

func (v *volume) update(volume int) {
	v.SetPercent(volume)
	go func() {
		app.Draw()
	}()
}
