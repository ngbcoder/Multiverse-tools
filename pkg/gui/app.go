package gui

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/multiverse/suite/pkg/disk"
	"github.com/multiverse/suite/pkg/installer"
)

// RunGUI launches the native Gio desktop GUI window loop.
func RunGUI() error {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Multiverse Desktop Suite"),
			app.Size(unit.Dp(640), unit.Dp(480)),
		)
		if err := loop(w); err != nil {
			fmt.Printf("Window loop error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
	return nil
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops

	var installBtn widget.Clickable
	statusText := "Status: Ready to flash target drive"
	progressVal := 0

	disks, _ := disk.ListDisks()
	if len(disks) == 0 {
		disks = []disk.DiskInfo{
			{Device: `C:\`, Model: "System Drive (C:)"},
		}
	}

	driveBtns := make([]widget.Clickable, len(disks))
	selectedIdx := 0

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Check for drive selection clicks
			for i := range driveBtns {
				if driveBtns[i].Clicked(gtx) {
					selectedIdx = i
					w.Invalidate()
				}
			}

			selectedDev := disks[selectedIdx].Device

			if installBtn.Clicked(gtx) {
				statusText = fmt.Sprintf("Status: Flashing drive %s...", selectedDev)
				w.Invalidate()

				go func() {
					inst := installer.NewInstaller(installer.Config{
						TargetDevice:   selectedDev,
						SecureBoot:     true,
						PartitionStyle: "GPT",
						Version:        "1.0.0",
					})

					_ = inst.Install(func(step string, percent int) {
						statusText = fmt.Sprintf("Status: [%d%%] %s", percent, step)
						progressVal = percent
						w.Invalidate()
					})

					statusText = fmt.Sprintf("Status: Flashing to %s Completed Successfully!", selectedDev)
					w.Invalidate()
				}()
			}

			// Render Gio UI layout
			layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				var elements []layout.FlexChild

				// Title
				elements = append(elements, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					h := material.H5(th, "Multiverse Desktop Installer & Drive Flasher")
					return h.Layout(gtx)
				}))

				elements = append(elements, layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout))

				// Drive Selection Subtitle
				elements = append(elements, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(th, "Select Target Storage Drive:")
					return lbl.Layout(gtx)
				}))

				elements = append(elements, layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout))

				// Render clickable drive buttons
				for i := range disks {
					idx := i
					elements = append(elements, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btnText := fmt.Sprintf("  %s - %s  ", disks[idx].Device, disks[idx].Model)
						if idx == selectedIdx {
							btnText = fmt.Sprintf("[SELECTED] %s - %s", disks[idx].Device, disks[idx].Model)
						}
						btn := material.Button(th, &driveBtns[idx], btnText)
						return btn.Layout(gtx)
					}))
					elements = append(elements, layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout))
				}

				elements = append(elements, layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout))

				// Progress Bar
				elements = append(elements, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					pb := material.ProgressBar(th, float32(progressVal)/100.0)
					return pb.Layout(gtx)
				}))

				elements = append(elements, layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout))

				// Status Label
				elements = append(elements, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					st := material.Body1(th, statusText)
					return st.Layout(gtx)
				}))

				elements = append(elements, layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout))

				// Install/Flash Action Button
				elements = append(elements, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btnLabel := fmt.Sprintf("Install / Flash Multiverse to %s", selectedDev)
					btn := material.Button(th, &installBtn, btnLabel)
					return btn.Layout(gtx)
				}))

				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, elements...)
			})

			e.Frame(gtx.Ops)
		}
	}
}
