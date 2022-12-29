package util

import (
	"encoding/json"
	"os"
	"proj2/png"
	"fmt"
	"log"
	"image"
)

type Task struct {
	InPath 	string 		`json:"inPath"`		//Specify the input directory
	OutPath string 		`json:"outPath"`	//Specify the output directory
 	Effects []string	`json:"effects"`	//Effects to apply for this task
}

// Reads all tasks from the json file at "../data/effects.txt" and returns a list
// of Task objects
func ParseTasks() []Task{
	// Load json file at "../data/effects.txt"
	effectsPathFile := fmt.Sprintf("../data/effects.txt")
	effectsFile, _ := os.Open(effectsPathFile)
	reader := json.NewDecoder(effectsFile)

	// Parse all of the tasks from the json file
	tasks := make([]Task, 0)
	for reader.More() {
		var task map[string]interface{}
		if err := reader.Decode(&task); err != nil {
			log.Println(err)
			return nil
		}
		effects := make([]string, 0)
		for _, effect := range (task["effects"]).([]interface{}){
			effects = append(effects, effect.(string))
		}
		tasks = append(tasks, Task{InPath: task["inPath"].(string), OutPath: task["outPath"].(string), Effects: effects})
	}
	return tasks
}

// Takes the resolutions directories to operate on, and a list of tasks execute
// generate a list of IamgeTask objects
func GenerateImageTasks(resolutions []string, tasks []Task) []*png.ImageTask{
	imageTasks := make([]*png.ImageTask, 0)
	for _, resolution := range resolutions {
		for _, task := range tasks {
			inPath := fmt.Sprintf("../data/in/%s/%s", resolution, task.InPath)
			outPath := fmt.Sprintf("../data/out/%s_%s", resolution, task.OutPath)
			// Load image from the input directory path
			pngImg, err := png.Load(inPath)
			if err != nil {
				panic(err)
			}
			// Save image, output path and list of effects to apply to the ImageTask object
			newImageTask := png.ImageTask{Image: pngImg, Effects: task.Effects, OutPath: outPath}
			newImageTask.Effects = task.Effects
			imageTasks = append(imageTasks, &newImageTask)
		}
	}
	return imageTasks
}

// Function to apply a single effect over an subBound of a given png Image
func ApplyEffect(pngImg *png.Image, effect string, bounds image.Rectangle){
	switch effect{
	case "S":
		pngImg.Sharpen(bounds)
	case "E":
		pngImg.EdgeDetection(bounds)
	case "B":
		pngImg.Blur(bounds)
	case "G":
		pngImg.Grayscale(bounds)
	default:
		log.Println("Unknown Effect")
		return
	}
}

// Function to set the sub-boundaries based on the number of threads and the id
// of a given thread
func SetBounds(id int, bound image.Rectangle, numThreads int) image.Rectangle {
	subBound := image.Rectangle{}
	numRows := bound.Max.Y - bound.Min.Y + 1
	fraction := numRows / numThreads
	subBound.Max.X = bound.Max.X
	subBound.Min.X = bound.Min.X
	subBound.Max.Y = bound.Max.Y - id*fraction
	if id == numThreads - 1 {
		subBound.Min.Y = bound.Min.Y
	} else {
		subBound.Min.Y = subBound.Max.Y - fraction
	}
	return subBound
}

