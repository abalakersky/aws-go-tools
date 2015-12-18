package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	logFile = flag.String("file", "yes", "Save output into file")
	t       = time.Now()
	dir, _  = filepath.Abs(filepath.Dir(os.Args[0]))
)

func main() {
	flag.Parse()

	fmt.Printf("WWIR", strconv{})

	name := dir + "/" + "output_" + strconv.FormatInt(t.Unix(), 10) + ".log"

	if *logFile == "yes" {
		f, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}
	w := bufio.NewWriter(f)

	for _, v := range my_slice {
		switch {
		case *logFile == "yes":
			fmt.Fprintln(w, v)
		case *logFile != "yes":
			fmt.Println(v)
		}
	}
	w.Flush()
}
