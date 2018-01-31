package openpomodoro

import (
	"fmt"
	"log"
	"os"
)

func debug(s string, i ...interface{}) {
	if os.Getenv("POMODORO_DEBUG") != "" {
		s = fmt.Sprintf("%s\n\n", s)
		log.Printf(s, i...)
	}
}
