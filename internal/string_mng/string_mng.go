package string_mng

import (
	"fmt"
	"os"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m\n" // light blue
	NoticeColor  = "\033[1;36m%s\033[0m\n" // green
	WarningColor = "\033[1;33m%s\033[0m\n" // yellow
	DebugColor   = "\033[0;36m[-debug-] %s\033[0m\n" // blue
	ErrorColor   = "\033[1;31m%s\033[0m\n" // red
)

func ColorizeErrorStart(){
	fmt.Println("\u001b[1;31m")
}

func ColorizeErrorEnd(){
	fmt.Println("\u001b[32m")
}

func PrintError(message string) {
	fmt.Printf(ErrorColor, message)
}

func PrintWarning(message string) {
	fmt.Printf(WarningColor, message)
}

func PrintNotice(message string) {
	fmt.Printf(NoticeColor, message)
}

func PrintNormal(message string) {
	fmt.Println(string("\u001b[0m"), message, string("\u001b[0m"))
}

func PrintDebug(message string) {
	fmt.Printf(DebugColor, message)
}

func PrintInfo(message string) {
	fmt.Printf(InfoColor, message)
}

func PrintBoolean(message bool){
	PrintNormal(fmt.Sprintf("%t : %v",message,message))
}
func UpdateLine(message string) {
	fmt.Printf("\r%s", message)
}

func PrintIntCarriage(length int) {
	fmt.Fprintf(os.Stderr, "\rQueue length: %d", length)
}

func PrintDynamicScanStatus(length int, url string, node int) {
	fmt.Fprintf(os.Stderr, "\rQueue length: %d    Analysing url: %s    found new %d nodes", length, url, node)
}

func ClearProgress() {
	fmt.Fprint(os.Stderr, resetTerminal())
}

func resetTerminal() string {
	return "\r\x1b[2K"
}

/*
func resetTerminal() string {
	return "\r\r"
}*/
