package main

import "fmt"

type Color string

const (
	// Bold
	BoldBlackFg   Color = "\u001b[30;1m"
	BoldRedFg     Color = "\u001b[31;1m"
	BoldGreenFg   Color = "\u001b[32;1m"
	BoldYellowFg  Color = "\u001b[33;1m"
	BoldBlueFg    Color = "\u001b[34;1m"
	BoldMagentaFg Color = "\u001b[35;1m"
	BoldCyanFg    Color = "\u001b[36;1m"
	BoldWhiteFg   Color = "\u001b[37;1m"
	// Regular
	BlackFg   Color = "\u001b[30m"
	RedFg     Color = "\u001b[31m"
	GreenFg   Color = "\u001b[32m"
	YellowFg  Color = "\u001b[33m"
	BlueFg    Color = "\u001b[34m"
	MagentaFg Color = "\u001b[35m"
	CyanFg    Color = "\u001b[36m"
	WhiteFg   Color = "\u001b[37m"
	// Reset
	Reset Color = "\u001b[0m"
)

func colorize(color Color, message string) string {
	return fmt.Sprint(string(color), message, string(Reset))
}

func debug(message string) {
	if *showDebug {
		fmt.Println(fmt.Sprintf("%s###%s %s", BoldBlueFg, Reset, message))
	}
}

func info(message string) {
	if !*noOutput {
		fmt.Println(fmt.Sprintf("%s>>>%s %s", BoldGreenFg, Reset, message))
	}
}

func warn(message string) {
	fmt.Println(fmt.Sprintf("%s***%s %s", BoldYellowFg, Reset, message))
}

func error(message string) {
	fmt.Println(fmt.Sprintf("%s!!!%s %s", BoldRedFg, Reset, message))
}
