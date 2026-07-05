package reqlogger

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const reset = "\033[0m"
const dot = "•"
const arrow = "»"

type color struct {
	r, g, b uint8
}

var Colors = map[string]color{ // Do not mutate
	"cyan":          {110, 190, 252},
	"green":         {99, 252, 76},
	"red":           {255, 66, 66},
	"deep_red":      {195, 6, 6},
	"magenta":       {181, 81, 178},
	"orange":        {255, 112, 51},
	"grey":          {125, 125, 125},
	"purple":        {142, 20, 145},
	"white":         {255, 255, 255},
	"greyish_white": {139, 148, 148},
	"yellow":        {255, 183, 48},
	"deep_blue":     {0, 55, 219},
	"light_blue":    {58, 151, 223},
	"pink":          {250, 50, 240},
	"black":         {0, 0, 0},
}

func colorText(c color, s string) string {
	return fmt.Sprintf(
		"\033[38;2;%v;%v;%vm%s%v",
		c.r, c.g, c.b,
		s, reset,
	)
}

func colorLabel(c color, s string) string {
	whiteColor := Colors["black"]
	return fmt.Sprintf(
		"\033[38;2;%d;%d;%d;48;2;%d;%d;%dm%s%v",
		whiteColor.r, whiteColor.g, whiteColor.b,
		c.r, c.g, c.b,
		s, reset,
	)
}

func RequestLog(ip string, method string, path string, statusCode int) {
	timeStr := colorText(Colors["grey"], time.Now().Format("03:04:05 PM"))
	methodStr := colorLabel(methodColor[method], padEqually(method))
	messageStr := `"` + path + `"` + " " + dot + " " + colorText(Colors["grey"], ip)
	println(timeStr + " " + dot + " " + methodStr + " " + dot + " " + colorText(getStatusCodeColor(statusCode), strconv.Itoa(statusCode)) + " " + arrow + " " + messageStr)
}
func padEqually(s string) string {
	// return " " + s + strings.Repeat(" ", 9-len(s)-1)
	t := 9 - len(s)
	r := (t + 1) / 2
	l := (t / 2)
	return strings.Repeat(" ", l) + s + strings.Repeat(" ", r)

}

func getStatusCodeColor(i int) color {
	if i < 200 {
		return Colors["greyish_white"]
	} else if i < 300 {
		return Colors["green"]
	} else if i < 400 {
		return Colors["yellow"]
	} else if i < 500 {
		return Colors["red"]
	} else if i >= 500 {
		return Colors["deep_red"]
	} else {
		return Colors["greyish_white"]
	}
}

var methodColor = map[string]color{
	"GET":     Colors["cyan"],
	"POST":    Colors["orange"],
	"PATCH":   Colors["yellow"],
	"PUT":     Colors["yellow"],
	"DELETE":  Colors["red"],
	"OPTIONS": Colors["greyish_white"],
	"HEAD":    Colors["greyish_white"],
}
