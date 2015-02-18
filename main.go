package main

import (
	"github.com/leafo/selfwatch/selfwatch"
)

func main() {
	recorder := selfwatch.NewRecorder()
	recorder.Bind()
}
