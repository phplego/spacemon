package util

import "github.com/fatih/color"

var gColorStyles = map[string]*color.Color{
	"title":    color.New(color.FgHiYellow, color.Bold),
	"dirs":     color.New(color.FgHiBlue),
	"diff":     color.New(color.FgHiMagenta),
	"free":     color.New(color.FgHiGreen, color.Bold),
	"pale":     color.New(color.FgHiYellow),
	"header":   color.New(color.FgBlue, color.Bold),
	"headerHi": color.New(color.FgHiBlue, color.Bold),
	"error":    color.New(color.FgHiRed),
}

func C(style, fmt string, a ...interface{}) string {
	colorStyle, _ := gColorStyles[style]
	return colorStyle.Sprintf(fmt, a...)
}
