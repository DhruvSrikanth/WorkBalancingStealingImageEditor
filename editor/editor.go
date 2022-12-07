package main

import (
	"fmt"
	"os"
	"proj3/scheduler"
	"strconv"
	"time"
)

const usage = "Usage: editor data_dir mode [number of threads]\n" +
	"data_dir = The data directory to use to load the images.\n" +
	"[mode]     = (ws) run the work stealing mode, (wb) run the work balancing mode\n" +
	"[number of threads] = Runs the parallel version of the program with the specified number of threads.\n" +
	"[threshold] = The threshold for the work stealing mode."

func main() {
	// Check for correct number of arguments
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}

	// Initialize the config
	config := scheduler.Config{DataDirs: "", Mode: "", ThreadCount: 0, Threshold: 0}
	config.DataDirs = os.Args[1]

	// Check for correct number of arguments
	if len(os.Args) > 5 {
		fmt.Println(usage)
		return
	}

	if len(os.Args) == 5 {
		// Work balancing mode
		config.Mode = os.Args[2]
		threads, _ := strconv.Atoi(os.Args[3])
		threshold, _ := strconv.Atoi(os.Args[4])
		config.ThreadCount = threads
		config.Threshold = threshold
	} else if len(os.Args) == 4 {
		// Work stealing mode
		config.Mode = os.Args[2]
		threads, _ := strconv.Atoi(os.Args[3])
		config.ThreadCount = threads
	} else {
		// Sequential mode
		config.Mode = "s"
	}

	// Time the image editor for all tasks in data/effects.txt
	start := time.Now()
	scheduler.Schedule(config)
	end := time.Since(start).Seconds()
	fmt.Printf("%.2f\n", end)

}
