package main

import (
	"os"

	"github.com/logmanager-oss/logveil/cmd/logveil"
)

func main() {
	os.Args = []string{"logveil", "-i", "/Users/Maciej/Work/Projects/logveil/tests/data/lm_export_test_input_raw_empty.csv", "-o", "output.txt"}
	logveil.Start()
}
