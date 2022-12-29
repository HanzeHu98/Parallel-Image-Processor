package scheduler

import(
	"proj2/png"
	"proj2/util"
	"fmt"
	"image"
	"strings"
)

// Run the image processing tasks using a pipeline scheduler (fan-out design)
func RunPipeline(config Config) {
	tasks := util.ParseTasks()
	resolutions := strings.Split(config.DataDirs, "+")
	numThreads := config.ThreadCount
	totalImages := len(tasks) * len(resolutions)

	// Create a generator channel and spawn a goroutine to generate imageTasks
	generator := make(chan *png.ImageTask)
	defer close(generator)
	go ImageTaskGenerator(generator, tasks, resolutions)

	// Create a aggregator channel for workers to push their results into
	aggregator := make(chan *png.ImageTask)
	defer close(aggregator)
	
	counter := make(chan int)
	defer close(counter)

	// Spawn numThreads number of pipeline workers 
	for i := 0; i < numThreads; i++{
		go PipelineWorker(generator, aggregator, numThreads)
	}

	// Spawn a thread to run the aggregator that saves their result, terminate
	// when all images results have been aggregated
	go Aggregator(aggregator, counter)
	aggregatedImage := 0
	for true{
		aggregatedImage = aggregatedImage + <-counter
		if aggregatedImage == totalImages{
			break
		}
	}
	return
}

// A Task Generator that loads tasks, generate imageTasks, and pushes them to the generator channel
func ImageTaskGenerator(generator chan<- *png.ImageTask, tasks []util.Task, resolutions []string) {
	for _, resolution := range resolutions {
		for _, task := range tasks {
			inPath := fmt.Sprintf("../data/in/%s/%s", resolution, task.InPath)
			outPath := fmt.Sprintf("../data/out/%s_%s", resolution, task.OutPath)

			pngImg, err := png.Load(inPath)
			if err != nil{
				panic(err)
			}
			imageTask := &png.ImageTask{Image: pngImg, Effects: task.Effects, OutPath: outPath}
			generator <- imageTask
		}
	}
}

// A PipelineWorker that calles the worker function when it receives an input from the generator channel
func PipelineWorker(generator <-chan *png.ImageTask ,aggregator chan<- *png.ImageTask, numThreads int){
	for true{
		task, more := <- generator
		if more{
			worker(task, aggregator, numThreads)
		} else{
			return
		}
	}
}

// A worker that spawns a set of mini-workers to handle a section of a given image
func worker(imageTask *png.ImageTask, aggregator chan<- *png.ImageTask, numThreads int){
	for index, effect := range imageTask.Effects {
		minicounter := make(chan int)
		for i := 0; i < numThreads; i++{
			go miniWorker(imageTask, effect, minicounter, util.SetBounds(i, imageTask.Image.Bounds, numThreads))
		}
		var finished = 0
		for true {
			finished = finished + <- minicounter
			if finished == numThreads {
				close(minicounter)
				break
			}
		}
		// If this is the last rendering subtask, then pipeline it to aggregator
		// Otherwise swap input and output image and go on to next effect
		if index == len(imageTask.Effects)-1 {
			aggregator <- imageTask
			return
		}
		imageTask.Image.SwapInOut()
	}
}

// A miniWorker that processes a single effect on a small area of the full image
func miniWorker(imageTask *png.ImageTask, effect string, counter chan<-int, subBounds image.Rectangle){
	util.ApplyEffect(imageTask.Image, effect, subBounds)
	counter <- 1
}

// A Result Aggregator that pulls from the aggregator channel, and saves the result to the output generate imageTasks, and pushes them to the generator channel
func Aggregator(aggregator <-chan *png.ImageTask, counter chan <- int){
	for true {
		result, more := <- aggregator
		if more {
			err := result.Image.Save(result.OutPath)
			if err != nil{
				panic(err)
			}
			counter <- 1
		} else {
			return
		}
	}
}