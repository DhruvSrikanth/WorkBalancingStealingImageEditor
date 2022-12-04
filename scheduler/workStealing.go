package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj3/concurrent"
	"proj3/png"
	"proj3/task"
	"strings"
)

// Run the work stealing model for generating and performing the tasks
func RunWorkStealing(config Config) {
	executor := concurrent.NewWorkStealingExecutor(config.ThreadCount, 1)

	dataDirs := strings.Split(config.DataDirs, "+")
	outputPath := "../data/out/%s_%s"
	inputPath := "../data/in/%s/%s"

	effectsPathFile := "../data/effects.txt"
	effectsFile, err := os.Open(effectsPathFile)
	if err != nil {
		panic(err)
	}
	defer effectsFile.Close()

	// Get the decoder
	reader := json.NewDecoder(effectsFile)

	// Decode the json requests in the effects file
	for {
		// Read the next request from the effects file
		// If there are no more requests, break
		job := Job{}
		err := reader.Decode(&job)
		if err != nil {
			break
		}

		// Process the task
		for _, dataDir := range dataDirs {
			inPath := fmt.Sprintf(inputPath, dataDir, job.InPath)
			outPath := fmt.Sprintf(outputPath, dataDir, job.OutPath)

			// Read the input file
			img, err := png.Load(inPath)
			if err != nil {
				panic(err)
			}

			// Add image task to the work pool
			imageTask := task.NewImageTask(img, outPath, job.Effects)
			executor.Submit(imageTask)
		}
	}

	// Shutdown the service
	executor.Shutdown()
}
