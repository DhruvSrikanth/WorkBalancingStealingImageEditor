package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj3/png"
	"proj3/task"
	"strings"
)

type Job = task.Job
type Image = png.Image

// Run the sequential model for generating and performing the tasks
func RunSequential(config Config) {
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

			// Create the image task
			imageTask := task.ImageTask{
				Image:      img,
				OutputPath: outPath,
				Effects:    job.Effects,
			}

			// Process the effects
			imageTask.ApplyEffects(img.Bounds.Min.Y, img.Bounds.Max.Y)

			// Save the output file
			imageTask.SaveResult()
		}
	}
}
