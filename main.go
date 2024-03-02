package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/mattn/go-isatty"
)

const (
	prgname = "decolor"
	prgver  = "1.0.1"
)

func printUsage() {
	fmt.Printf(prgname + " v" + prgver + "\n" +
		"Text decolorizer\n" +
		"Usage: " + prgname + " [options]\n" +
		"  |piped input|      Piped text is decolorized\n" +
		"  FILENAME           Decolorize given file\n" +
		"  -?, -h, --help     Print this usage page\n")
	os.Exit(0)
}

func isGitBashOnWindows() bool {
	return runtime.GOOS == "windows" && strings.HasPrefix(os.Getenv("MSYSTEM"), "MINGW")
}

func hasPipedInput() bool {
	stat, _ := os.Stdin.Stat() // Check if anything was piped in
	if isGitBashOnWindows() {
		// Git Bash on Windows handles input redirection differently than other shells. When a program
		// is run without any input or arguments, it still treats the input as if it were piped from an
		// empty stream, causing the program to consider it as piped input and hang. This works around that.
		if !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
			return true
		}
	} else {
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			return true
		}
	}
	return false
}

func loadAndDecolorize(filename string) {
	// Read content from the given file
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Remove color codes from file content and print
	decolorizedText := color.ClearCode(string(fileBytes))
	fmt.Printf(decolorizedText)
}

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-?", "-h", "--help":
			printUsage()
		default:
			loadAndDecolorize(os.Args[1])
		}
	} else if hasPipedInput() {
		// Process piped input
		rawBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
		}

		// Remove color escape codes in piped input, then print
		decolorizedText := color.ClearCode(string(rawBytes))
		fmt.Printf(decolorizedText)
	} else {
		printUsage()
	}
}
