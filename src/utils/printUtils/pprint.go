package printUtils

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	No_color = "30"
	Red      = "31"
	Green    = "32"
	Yellow   = "33"
	Blue     = "34"
	Magenta  = "35"
	Cyan     = "36"
	Grey     = "37"
	Black    = "38"

	No_bg      = "40"
	On_red     = "41"
	On_green   = "42"
	On_yellow  = "43"
	On_blue    = "44"
	On_magenta = "45"
	On_cyan    = "46"
	On_grey    = "47"
	On_black   = "48"

	Bold = "1"
)

// returns a colorized version of s
// if it fails to read the color, returns the initial string
func Color(s, color string) string {
	colors := strings.Split(strings.ToLower(color), " ")
	var code string   = "\033["
	var first_el bool = true

	for _, c := range colors {
		if ! first_el {
			code += ";"
		} else {
			first_el = false
		}

		switch c {
		case "no_color":
			code += No_color
		case "red":
			code += Red
		case "green":
			code += Green
		case "yellow":
			code += Yellow
		case "blue":
			code += Blue
		case "magenta":
			code += Magenta
		case "cyan":
			code += Cyan
		case "grey":
			code += Grey
		case "black":
			code += Black
		case "no_bg":
			code += No_bg
		case "on_red":
			code += On_red
		case "on_green":
			code += On_green
		case "on_yellow":
			code += On_yellow
		case "on_blue":
			code += On_blue
		case "on_magenta":
			code += On_magenta
		case "on_cyan":
			code += On_cyan
		case "on_grey":
			code += On_grey
		case "on_black":
			code += On_black
		case "bold":
			code += Bold
		default:
			return s
		}
	}

	return code + "m" + s + "\033[0m"
}

// Tries to get the dimensions of the terminal. If it fails, then it returns 64x32
func getTerminalDimensions() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 64, 32
	} else {
		arr := strings.Split(string(out), " ")
		if len(arr) < 2 {
			return 128, 32
		}
		rows, err1 := strconv.Atoi(arr[0])
		cols, err2 := strconv.Atoi(arr[1][:len(arr[1])-1])
		if err1 != nil || err2 != nil {
			return 92, 32
		}
		return cols, rows
	}
}

func PrintCentered(thing string) {
	cols, _ := getTerminalDimensions()
	if len(thing) < cols {
		var beforeSpace int = (cols - len(thing)) / 2
		fmt.Printf("\033[%vG%v\n", beforeSpace, thing)
	} else {
		arr := strings.Split(thing, " ")
		var line string
	accumulator:
		for _, subS := range arr {
			if len(line) > 0 && len(line)+len(subS) > cols { // the line is filled
				PrintCentered(line)  // we print
				line = subS          // we renew
				continue accumulator // we go to the next loop
			}
			line += subS           // we add
			for len(line) > cols { // we overfill
				fmt.Println(line[:cols]) // we print what fits
				line = line[cols:]       // we delete what we printed
			}
			line += " " // we add the space
		}
		if len(line) > 1 { // we still have something to print, and we know its shorter than [cols]
			PrintCentered(line[:len(line)-1])
		}
	}
}

func PrintHLine(char rune) {
	cols, _ := getTerminalDimensions()
	c := string(char)
	for i := 0; i < cols; i++ {
		fmt.Printf(c)
	}
	fmt.Printf("\n")
}
