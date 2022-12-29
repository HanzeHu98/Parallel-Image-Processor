package scheduler

import (
	"proj2/png"
	"proj2/util"
	"sync"
	"strings"
)

type bspWorkerContext struct {
	// Define the necessary fields for your implementation
	taskCounter    	int
	effectCounter 	int
	totalTasks     	int
	imageTasks   	[]*png.ImageTask
	NumThreads   	int
	mt           	sync.Mutex
	cond         	*sync.Cond
	finishedCounter int
}

func NewBSPContext(config Config) *bspWorkerContext {
	//Initialize the context
	ctx := bspWorkerContext{taskCounter: 0, effectCounter: 0, finishedCounter: 0}
	ctx.NumThreads = config.ThreadCount

	resolutions := strings.Split(config.DataDirs, "+")
	tasks := util.ParseTasks()
	ctx.imageTasks = util.GenerateImageTasks(resolutions, tasks)
	ctx.totalTasks = len(tasks) * len(resolutions)

	ctx.cond = sync.NewCond(&ctx.mt)
	return &ctx
}

// This implementation of the builk synchronous parallel model has each
// goroutine working on a subsection of one effect, and when an effect is 
// applied, synchronize by checking if the task is done, if done, save to output,
// otherwise switch input and output iamges and continue on to the next effect/superstep.
func RunBSPWorker(id int, ctx *bspWorkerContext) {
	for {
		currTask := ctx.imageTasks[ctx.taskCounter]
		currSubTask := currTask.Effects[ctx.effectCounter]
		subBound := util.SetBounds(id, currTask.Image.Bounds, ctx.NumThreads)
		util.ApplyEffect(currTask.Image, currSubTask, subBound)

		// Syncrhonization step
		ctx.mt.Lock()
		ctx.finishedCounter++
		if ctx.finishedCounter == ctx.NumThreads {
			// Perform synchronization, if all effects done for this task, save,
			// otherwise swap input and output image.
			if ctx.effectCounter == len(currTask.Effects)-1 {
				err := currTask.Image.Save(currTask.OutPath)
				if err != nil {
					panic(err)
				}
				ctx.taskCounter++
				ctx.effectCounter = 0
			} else {
				currTask.Image.SwapInOut()
				ctx.effectCounter++
			}
			ctx.finishedCounter = 0
			ctx.cond.Broadcast()
		} else {
			ctx.cond.Wait()
		}
		ctx.mt.Unlock()

		// If all tasks have been completed, exit by breaking the for loop
		if ctx.taskCounter >= ctx.totalTasks {
			break
		}
	}
}
