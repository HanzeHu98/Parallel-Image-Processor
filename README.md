# Parallel Image Processing

## Introduction
This is an image processing system programmed in Golang, applies a range of effects using convolution.
See word document in benchmark subdirectory for more details on the implementation and analysis.

## How to run
Images to be processed and effects to apply should be included in the `data/effects.txt` document, in the following format:
```
{"inPath": "Image Name", "outPath": "Output Image Name", "effects": ["G","E","S"]}
```
Where the applicable effects are:
- "G": Grayscale Effect
- "E": Edge Detection Effect
- "S": Sharpening Effect
- "B": Blur Effect

Images should be placed inside the small, mixture or large subdirectories based on their sizes/resolution. You could specify which directory to read from when running the program. <br>
To run the program, use the following command:<br>
```go run proj2/editor small bsp 2```
- small can be replacd by **small**, **mixture** or **large**, based on the data folder you wish to read from
- bsp can be replaced by **bsp**, **pipeline** or **sequential**, based on the parallization technique you wish to use
- 2 can be replaced by any number and is used to indicate the number of threads to use

You can also run the benchmark program by using the following command: <br>
`python3 benchmark/plot.py`, this will produce two speedup graphs as below, note that this could take a while to finish running. 

## Description
The overall project implements a Image Procesing System using Golang. <br>
The png folder contains the definitions for working with images in golang, including loading, saving and updating png images, as well as the code for applying affects using convolutions to individual images. All code in the png folder are sequential.<br>
The editor folder contains editor.go, which is used to parse in configuration data such as type of scheduler scheme to use and number of cores to use if it the scheme can run in parallel.<br>
In the scheduler folder, there are three different implementations for three different types of scheduling schemes: pipeline- which implements a Fan-Out Fan-In design using channels for creation of iamge processing tasks and task result aggregation; bulk synchronous parallel (bsp)- which implements a bsp model technique to process the images step by step and synchronize between steps; and sequential- which processes all input images and effects sequentially, one by one.<br>

## Speedup Graph and Analysis
#### Pipeline Implementation
As per the graph below, we can see that the general speed-up trend appears to decrease as we increase the number of threads on the local machine. I think this trend may be due to the fact that since we are spawning numThreads number of goroutines as workers, who are also spawning numThreads number of sub-goroutines, as the number of threads increases, the amount of overhead also increases drastically. The randomness that’s apparent on the Linux Cluster may also be due to having too many goroutines to manage and more randomness in the thread management of the cluster.<br>
![Image](https://github.com/mpcs52060-aut22/project-2-HanzeHu98/blob/main/proj2/benchmark/speedup-pipeline.png)

#### Bulk Synchronous Parallelization Implementation
The graphs below, one ran on the linux cluster and the other on my local machine (m2 macbook air), shows that for all input image sizes, the bsp scheduling scheme is effective in reducing the runtime of overall processing, and can boost the performance of large datasets by up to 2 - 3 times. The dip in speedup for 8 threads on the linux cluster is likely an outlier, as it only appeared for small and mixture input sizes, and the trend did not show up on my local machine. While on the local machine, since it only had 8 cores, speedup was not significantly increased after 8 cores, and even slowed down slightly due to the increased overhead.<br>
![image](https://github.com/mpcs52060-aut22/project-2-HanzeHu98/blob/main/proj2/benchmark/speedup-bsp.png)

## Analysis
#### Bottlenecks and HotSpot
The hotspots in my sequential program mainly comes from applying the convolution filters to individual image pixels, due to the vast number of pixels in the input images (especially large images), the convolution and image processing step is expected to draw the most amount of computation resources. The bottleneck of my sequential program would come from reading and writing files. Since sequential implementation meant that no more than 1 image could be processed

#### Performance Comparison
In general, the pipeline implementation performed better than the bsp implementation, however, as we increase the number of threads, the speedup improves with the bsp implementation whereas the speedup decreases with the pipeline implementation, to the point where at 12 cores, bsp is performing better than pipeline (or at least at the same level).

#### Potential Improvements
Under the BSP scheduling scheme, we can potentially improve performance further would be to use separate go routines that are dedicated to reading and writing files, this way, the BSP could potentially start processing images before all images are loaded initially, and can save images as they are processed. 
