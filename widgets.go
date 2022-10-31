package main

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Field struct {
	widget.Editor

	Invalid bool
	old     string
}

func (ed *Field) Changed() bool {
	newText := ed.Editor.Text()
	changed := newText != ed.old
	ed.old = newText

	return changed
}

func (ed *Field) SetText(s string) {
	ed.Invalid = false
	ed.old = s
	ed.Editor.SetText(s)
}

func (ed *Field) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	borderWidth := float32(0.5)
	if ed.Editor.Focused() {
		borderWidth = 2
	}

	borderColor := color.NRGBA{A: 107}

	if ed.Invalid && ed.Text() != "" {
		borderColor = color.NRGBA{R: 200, A: 0xFF}
	} else {
		if ed.Editor.Focused() {
			borderColor = th.Palette.ContrastBg
		}
	}

	return widget.Border{
		Color:        borderColor,
		CornerRadius: unit.Dp(4),
		Width:        unit.Dp(borderWidth),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(4)).Layout(gtx,
			material.Editor(th, &ed.Editor, "").Layout)
	})
}

func Heading(th *material.Theme, txt string) material.LabelStyle {
	label := material.Label(th, th.TextSize*18.0/16.0, txt)
	label.Font.Weight = text.Bold

	return label
}

func Subheading(th *material.Theme, txt string) material.LabelStyle {
	label := material.Label(th, th.TextSize, txt)
	label.Font.Weight = text.Medium

	return label
}
