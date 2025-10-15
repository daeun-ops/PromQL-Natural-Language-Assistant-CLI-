package utils

import (
	"fmt"
)

var (
	bold   = "\033[1m"
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

func Section(title string) {
	fmt.Printf("\n%s== %s ==%s\n", bold, title, reset)
}

func Ok(msg string) {
	fmt.Printf("%s✔ %s%s\n", green, msg, reset)
}

func Warn(msg string) {
	fmt.Printf("%s! %s%s\n", yellow, msg, reset)
}

func Fail(msg string) {
	fmt.Printf("%s✖ %s%s\n", red, msg, reset)
}

