package main

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/amarillier/KrankyBearTimer/internal/startup"
)

// windowsRenderingRow is the About window's Windows-only rendering line and Mesa3D credit: which OpenGL path this
// launch uses (hardware, or the bundled Mesa3D software fallback and its rasterizer). Empty on other platforms.
func windowsRenderingRow() fyne.CanvasObject {
	if runtime.GOOS != "windows" {
		return container.NewVBox()
	}
	credit := widget.NewLabel("The Windows package includes Mesa 3D (MIT and other licences) as a software OpenGL fallback: see mesa-fallback\\MESA-NOTICE.txt.")
	credit.Wrapping = fyne.TextWrapWord
	credit.Alignment = fyne.TextAlignCenter
	return container.NewVBox(container.NewCenter(widget.NewLabel(renderingStatusText())), credit)
}

// aboutWindowSize leaves room for the extra Windows-only rendering rows.
func aboutWindowSize() fyne.Size {
	if runtime.GOOS == "windows" {
		return fyne.NewSize(500, 400)
	}
	return fyne.NewSize(500, 300)
}

// renderingStatusText reports which OpenGL path this launch uses: real hardware, or the bundled Mesa3D software fallback with
// the rasterizer (llvmpipe fast / softpipe safe) and the pinned Mesa version when known.
func renderingStatusText() string {
	active, version, driver := startup.MesaFallbackActive()
	if !active {
		return "Rendering: hardware OpenGL"
	}
	if version == "" {
		return "Rendering: Mesa3D software OpenGL (" + driver + ")"
	}
	return "Rendering: Mesa3D software OpenGL (" + driver + ", Mesa " + version + ")"
}
