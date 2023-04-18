package util

import "github.com/fatih/color"

var gColorStyles = map[string]*color.Color{
	"pale":     color.New(color.FgHiYellow),
	"header":   color.New(color.FgBlue).Add(color.Bold),
	"headerHi": color.New(color.FgHiBlue).Add(color.Bold),
}

func C(style, fmt string, a ...interface{}) string {
	colorStyle, _ := gColorStyles[style]
	return colorStyle.Sprintf(fmt, a...)
}
