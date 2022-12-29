package scheduler

import (
	"proj2/util"
	"strings"
	"proj2/png"
)

// Run the image processing tasks sequentially
func RunSequential(config Config) {
	// Parse all tasks that requires to be performed, as well as the resolution
	// level to perform them at
	tasks := util.ParseTasks()
	resolutions := strings.Split(strings.TrimSpace(config.DataDirs), "+")

	// Generate a list of imageTasks to execute and iterate over them to work on each
	imageTasks := util.GenerateImageTasks(resolutions, tasks)
	for _, task := range imageTasks{
		sequentialWorker(task)
	}
	return
}

// A sequential worker that applies all effects in an ImageTask one by one until
// all effects have been applied, it then saves the task to the output path
func sequentialWorker(imageTask *png.ImageTask){
	for index, effect := range imageTask.Effects{
		util.ApplyEffect(imageTask.Image, effect, imageTask.Image.Bounds)
		if index == len(imageTask.Effects) - 1{
			if err := imageTask.Image.Save(imageTask.OutPath); err != nil{
				panic(err)
			}
			return
		}
		imageTask.Image.SwapInOut()
	}
}
