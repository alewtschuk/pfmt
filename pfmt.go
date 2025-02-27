package pfmt

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alewtschuk/dsutils"
)

var colorMap = make(map[int]string, 464)
var fontMap = make(map[string]string, 5)
var keys []int = initColorMap()

// Initialize the color map with ANSI escape codes
// returns a slice of the keys in order of insertion
// function can be called without return value being assigned to a variable
// and will still initialize the color map
func initColorMap() []int {
	var keys []int // Slice to hold the keys in order of insertion
	for i := -231; i <= 231; i++ {
		keys = append(keys, i)
		if i == -1 {
			colorMap[i] = "\x1b[0m" // Add reset color as key -1
		} else if i < 0 {
			colorMap[i] = fmt.Sprintf("\x1b[48;5;%dm", -i) // ANSI escape character for background coloring with 256 options. Note that greyscale 256 values are not supported
		} else if i == 0 {
			colorMap[i] = "\x1b[38;5;0m" // ANSI escape character for resetting color
		} else {
			colorMap[i] = fmt.Sprintf("\x1b[38;5;%dm", i) // ANSI escape character for coloring with 256 options.
		}
	}
	//fmt.Println("Map initialized")
	return keys
}

// Initialize the font map with ANSI escape codes
// for font styling
func InitFontMap() {
	fontMap["RESET"] = "\x1b[0m"
	fontMap["BOLD"] = "\x1b[1m"
	fontMap["UNDERLINE"] = "\x1b[4m"
	fontMap["ITALIC"] = "\x1b[3m"
}

// Printc prints the formatted string in the specified color.
// The color will reset after printing the specified string
func Printc(format string, color int) (int, error) {
	var newline string = ""
	//Checks if the format string ends with a newline to handle newline edgecase
	if strings.HasSuffix(format, "\n") {
		format = strings.TrimSuffix(format, "\n")
		newline = "\n"
	}

	var s string = colorMap[color] + format + colorMap[-1] + newline // Containerizes as its own string
	return fmt.Print(s)
}

// Printcln prints the formatted string in the specified color with a new line.
// The color will reset after printing the specified string
func Printcln(format string, color int) (int, error) {
	var s string = colorMap[color] + format + colorMap[-1]
	return fmt.Println(s)
}

// Printm prints the formatted string in the specified colors.
// Utilizes the format string "%h"(hue) to determine where to apply the colors.
func Printmc(format string, a ...int) {
	var colorCount int = 0
	var colorStrings []string
	var newline string = ""
	//Checks if the format string ends with a newline to handle newline edgecase
	if strings.HasSuffix(format, "\n") {
		format = strings.TrimSuffix(format, "\n")
		newline = "\n"
	}
	colorStrings = strings.Split(format, "%h")      // Split the format string into substrings at "%h"
	colorStrings = dsutils.Remove("", colorStrings) // Remove any empty strings from the slice\

	// Iterate over the substrings, color them and print them
	for _, substring := range colorStrings {
		if colorCount < len(a) {
			if colorCode, exists := colorMap[a[colorCount]]; exists { //If the color exists print the color, substring, and reset
				fmt.Print(colorCode)
				fmt.Print(substring)
				fmt.Print(colorMap[-1])
				colorCount++
				continue
			}
		}
	}
	fmt.Print(colorMap[-1] + newline)
}

// Printmcln utilizes the format string "%h"(hue) to determine where to apply colors to the string.
// Prints the formatted string in the specified colors and adds a new line.
func Printmcln(format string, a ...int) {
	var colorCount int = 0
	var colorStrings []string
	colorStrings = strings.Split(format, "%h")
	colorStrings = dsutils.Remove("", colorStrings)

	for _, substring := range colorStrings {
		if colorCount < len(a) {
			if colorCode, exists := colorMap[a[colorCount]]; exists {
				fmt.Print(colorCode)
				fmt.Print(substring)
				fmt.Print(colorMap[-1])
				colorCount++
				continue
			}
		}
	}
	fmt.Print(colorMap[-1])
	fmt.Println() // This is literally the only change from Printmc
}

// Prints the string with specified foreground and background colors.
// The foreground must be positive
// The background must be negative
func Printcfb(format string, foreground int, background int) (int, error) {
	var s string
	var newline string = ""

	if strings.HasSuffix(format, "\n") {
		format = strings.TrimSuffix(format, "\n")
		newline = "\n"
	}

	if foreground > 0 && background < -1 {
		s = fmt.Sprintf("%s%s%s%s", colorMap[foreground], colorMap[background], format, colorMap[-1])
	} else {
		return 0, fmt.Errorf("invalid color combination")
	}

	return fmt.Print(s + newline)
}

// Prints the string with specified foreground and background colors.
// The foreground must be positive
// The background must be negative
func Printcfbln(format string, foreground int, background int) (int, error) {
	var s string
	var newline string = ""

	if strings.HasSuffix(format, "\n") {
		format = strings.TrimSuffix(format, "\n")
		newline = "\n"
	}

	if foreground > 0 && background < -1 {
		s = fmt.Sprintf("%s%s%s%s", colorMap[foreground], colorMap[background], format, colorMap[-1])
	} else {
		return 0, fmt.Errorf("invalid color combination")
	}

	return fmt.Println(s + newline)
}

// ApplyColor applies the specified color to the input string.
// Returns the string with the specified color applied.
func ApplyColor(format string, color int) string {
	format = StripColor(format)
	format = colorMap[color] + format + colorMap[-1]
	return format
}

// Removes the color from the input string
func StripColor(format string) string {
	for k := keys[0]; k <= keys[len(keys)-1]; k++ {
		if k != -1 { // Don't remove the reset sequence itself
			format = strings.Replace(format, colorMap[k], "", -1) // Also remove background colors
		}
	}
	return format
}

// Checks if the input color is valid
func IsColorValid(color int) bool {
	if _, exists := colorMap[color]; exists {
		return true
	} else {
		return false
	}
}

// Displays all available colors integers -232-231 in their corresponding color
func AvailableColors() {
	for k := keys[0]; k <= keys[len(keys)-1]; k++ {
		fmt.Print(colorMap[k] + strconv.Itoa(k) + colorMap[-1])
		if k >= -1 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}

// Loader function to display a loading animation, just used as a temp
// function till possible later integration
// func Loader() {
// 	for i := 0; i < 10; i++ {
// 		fmt.Print("\x1b[1m" + strings.Repeat("\r⠋", i))
// 		time.Sleep(30 * time.Millisecond)
// 		fmt.Print("\x1b[1m" + strings.Repeat("\r⠙", i))
// 		time.Sleep(30 * time.Millisecond)
// 		fmt.Print("\x1b[1m" + strings.Repeat("\r⠹", i))
// 		time.Sleep(30 * time.Millisecond)
// 		fmt.Print("\x1b[1m" + strings.Repeat("\r⠸", i))
// 		time.Sleep(30 * time.Millisecond)
// 		fmt.Print("\x1b[1m" + strings.Repeat("\r⠼", i))
// 		time.Sleep(30 * time.Millisecond)
// 		fmt.Print("\x1b[1m" + strings.Repeat("\r⠴", i))
// 		time.Sleep(30 * time.Millisecond)
// 		fmt.Print("\x1b[1m" + strings.Repeat("\r⠦", i))
// 		time.Sleep(30 * time.Millisecond)
// 		fmt.Print("\x1b[1m" + strings.Repeat("\r⠧", i))
// 		time.Sleep(30 * time.Millisecond)
// 		fmt.Print("\x1b[1m" + strings.Repeat("\r⠇", i))
// 		time.Sleep(30 * time.Millisecond)
// 		fmt.Print("\x1b[1m" + strings.Repeat("\r⠏", i))
// 		time.Sleep(30 * time.Millisecond)
// 	}
// 	fmt.Print("\r ")
// 	fmt.Println()
// }
