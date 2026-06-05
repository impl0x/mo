package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const reset = "\033[0m"
const dot = "•"
const arrow = "»"

type color struct {
	r, g, b uint8
}

var Colors = map[string]color{
	"blue":       {0, 0, 255},
	"cyan":       {50, 125, 255},
	"green":      {21, 140, 10},
	"red":        {175, 0, 0},
	"magenta":    {181, 81, 178},
	"orange":     {255, 123, 79},
	"grey":       {125, 125, 125},
	"purple":     {142, 20, 145},
	"white":      {255, 255, 255},
	"yellow":     {194, 158, 0},
	"deep_blue":  {0, 55, 219},
	"light_blue": {58, 151, 223},
	"pink":       {250, 50, 240},
}

//? Levels

type Level struct {
	Color color
	Label string
}

var info = Level{
	Color: Colors["light_blue"],
	Label: "info",
}
var warn = Level{
	Color: Colors["yellow"],
	Label: "warn",
}

var success = Level{
	Color: Colors["green"],
	Label: "success",
}

var error = Level{
	Color: Colors["red"],
	Label: "error",
}

var fatal = Level{
	Color: Colors["purple"],
	Label: "fatal",
}

// ? extra helper funcs

func log(lvl Level , v []string){
	if len(v)==0{
		_loggerError("Not enough values passed in logger",false)
		return
	}
	msg := v[0]
	var kwargs []string =nil
	if len(v) > 1 {
		kwargs=v[1:]
	}
	println(lvl.makeLogMessage(msg, kwargs))
}

func _loggerError(msg string, fatal bool) {
	logFatal := Level{
		Color: Colors["purple"],
		Label: "logger",
	}
	println(logFatal.makeLogMessage(msg, nil))
	if fatal {
		os.Exit(1)
	}
}


// ? Logging internals

func colorText(c color, s *string) {
	*s = fmt.Sprintf(
		"\033[38;2;%v;%v;%vm%s%v",
		c.r, c.g, c.b,
		*s, reset,
	)
}

func colorLabel(c color, s *string) {
	whiteColor := Colors["white"]
	*s = fmt.Sprintf(
		"\033[38;2;%d;%d;%d;48;2;%d;%d;%dm%s%v",
		whiteColor.r, whiteColor.g, whiteColor.b,
		c.r, c.g, c.b,
		*s, reset,
	)

}
func (lvl Level) makeLogMessage(message string, args []string) string {
	nowTime := time.Now()
	timeString := nowTime.Format("03:04:05 PM")
	colorText(Colors["grey"],&timeString)

	// ? Logic for having a evenly spaced background
	t := 9 - len(lvl.Label)
	r := (t + 1) / 2
	l := (t / 2)
	lbStr := strings.Repeat(" ", l) + lvl.Label + strings.Repeat(" ", r)

	colorLabel(lvl.Color, &lbStr)
	
	finalString := fmt.Sprintf("%v %v %v %v %v",
		timeString, dot, lbStr, arrow, message,
	)

	
	if len(args) != 0 {
		var toAddBuilder strings.Builder
		toAddBuilder.Write([]byte(" " + dot))

		var kColor color
		for i := 0; i < len(args); i += 2 {
			k, v := args[i], args[i+1]

			kColor = Colors["grey"]

			colorText(kColor, &k)

			toAddBuilder.Write([]byte(" " + k + ": " + v))
			if i < len(args)-2 {
				toAddBuilder.Write([]byte(","))
			}
		}
		finalString += toAddBuilder.String()
	}
	return finalString
}

//? Func definations for levels
func Info(v ...string) {
	log(info,v)
}
func Warn(v ...string) {
	log(warn,v)
}
func Success(v ...string) {
	log(success,v)
}
func Error(v ...string) {
	log(error,v)
}
func Fatal(v ...string) {
	log(fatal,v)
	os.Exit(1)
}

func Custom(lvl Level, v ...string){
	log(lvl,v)
}

// * usage guide:
// do logger.[level] after importing the logger package
// example, logger.info("Message","key","value")
// key, value args are optional and are displayed as key: value.
// for custom function, check the colors from {Colors}