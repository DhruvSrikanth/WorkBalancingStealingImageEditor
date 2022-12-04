package task

import (
	"proj3/png"
	"sync"
)

type Image = png.Image

type Job struct {
	InPath  string   `json:"inPath"`
	OutPath string   `json:"outPath"`
	Effects []string `json:"effects"`
}

type ImageTask struct {
	Image      *Image
	OutputPath string
	Effects    []string
	cond       *sync.Cond
	done       bool
}

// Create a new task
func NewImageTask(image *Image, outputPath string, effects []string) interface{} {
	// Create a new condition variable
	cond := sync.NewCond(&sync.Mutex{})
	// Create a new task
	task := &ImageTask{
		Image:      image,
		OutputPath: outputPath,
		Effects:    effects,
		cond:       cond,
		done:       false,
	}
	return task
}

// Apply the effects to the image
func (t *ImageTask) ApplyEffects(startY, endY int) {
	// Process all effects and swap the buffers after each effect
	for _, effect := range t.Effects {
		t.Image.ApplyEffect(effect, startY, endY)
		t.Image.Swap()
	}
	// Swap the buffers back to the output
	t.Image.Swap()
}

func (t *ImageTask) SaveResult() {
	// Save the output file
	err := t.Image.Save(t.OutputPath)
	if err != nil {
		panic(err)
	}
}

// Run the task
func (task *ImageTask) Run() {
	// Process the effects
	task.ApplyEffects(0, task.Image.Bounds.Max.Y)
	// Save the output file
	task.SaveResult()
	// Signal that the task is complete
	task.Done()
}

// Wait for the task to complete
func (f *ImageTask) Get() interface{} {
	// Wait for the barrier to be signaled
	f.cond.L.Lock()
	if !f.done {
		f.cond.Wait()
	}
	f.cond.L.Unlock()
	return nil
}

// Indicates that the task is complete
func (f *ImageTask) Done() {
	// Signal that the barrier is complete
	f.cond.L.Lock()
	f.done = true
	f.cond.Signal()
	f.cond.L.Unlock()
}
