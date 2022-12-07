# Project: Image Editor

**Important Note 1** - 

In the project specifications, it is mentioned that all files provided to us are within the `proj3` directory. However, when pulling the repository, there was no `proj3` directory. Since I didn't know if this was on purpose, I have continued without a `proj3` directory. All paths have been appropriately changed and the project for the advanced feature (MapReduce) is inside the `conjugate_gradient` directory. I would have preferred to organize the directory structure a little differently (divide the two projects into two distinct folders) however, I was not sure of the requirements so I have made as little changes to the initial directory structure provided as possible.

**Important Note 2** -

The advanced feature utilizing MapReduce can be found within the `conjugate_gradient` directory. This is a completely separate application. The report for this can be found within the **Advanced Feature** section of this README. I'd like to thank Umang Bhatia for brainstorming ideas with me for the advanced feature project. The conjugate gradient sparse matrix solver was a project I previously did in `C++`, `OpenMP` and `MPI`. The motivation behind reimplementing this in `Go` was to explore how `MPI_Allreduce`, which I used in my previous implementation, can be implemented in `Go`.


### Project Description and Motivation - 

This project is an implementation of an image editor. The motivation behind this project is to improve the compute time taken by image editors. In a digital world, appearances online have become so important, from Zoom meeting filters to liven a call, to Instagram pages to promote products. With the increase in the ubiquity of multi-core machines, parallel computing offers a new paradigm of programming that can take advantage of more compute power to drastically decrease the wait time to produce images that may just brighten someone's day.

In this image editor, images are transformed based on a sequence of effects that are to be applied to each image along with the name of the image to be transformed and the name of the transformed image to be saved. These transformations (effects) are carried by applying specific kernels through a convolution operation over the entire image. For all effects, a 3x3 kernel is convolved over patches (iterating over the entire image). 

The editor processes image tasks specified in three different ways i.e. sequential, work stealing (ws) and work balancing (wb) parallel mode. In work balancing, a thread determines whether it needs to balance work with another thread chosen at random. The decision to balance is based on the balance threshold specified as a command-line argument. If the difference in the number of tasks assigned to a thread's queue and another queue chosen at random, is greater than or equal to the balancing threshold, then the thread will balance work with the other chosen thread. In the work stealing algorithm, upon the completion of all tasks assigned to a thread, the thread will steal work from another thread chosen at random (if of course the other thread's queue contains image tasks to be processed).

### Important System Components - 

1. Data (Input):
    - **I/O takes places in the `data` directory. If there is no `data` directory, you can create one using `mkdir data` while inside the project directory.**
    - Each editor task is provided within the directory `data/effects.txt` file. An example of an editor task is shown below (multiple tasks can be specified by placing one on each line of the `effects.txt` file and are in the `JSON format` - 

        ```JSON
        {"inPath": "IMG_2029.png", "outPath": "IMG_2029_Out.png", "effects": ["G","E","S","B","M"]}
        ``` 
    
    - Input images must be provided within the a directory inside `data/in` directory. If one does not exist, you can create it from within `data` using `mkdir in`. Similarly for the directory. For example, to process an image, we can create  `data/in/small` and place the image inside this directory.  

    - Input images are provided in the `PNG` format.

    - The data can be found [here](https://drive.google.com/drive/folders/17A9GdVSkzKJtFFu2MNuozja8soqulqkD?usp=sharing). **Make sure to copy the `data` folder into the root directory of this project.**


2. Data (Output):

    - **I/O takes places in the `data` directory. If there is no `data` directory, you can create one using `mkdir data` while inside the project directory.**

    - After processing images, the transformed images are saved within the `data/out` directory. If one does not exist, you can create it from from within `data` using `mkdir out`. 

    - Processed images are saved in the `PNG` format.

3. Effects:
    - The following effects can be applied (the appropriate identifier to be used is specified in parenthesis for each effect) - 

        - Grayscale (G) - Each pixel's color channels are computed by averaging over the each the original pixel's color channel values.

        - Sharpen (S) - The following kernel is used as part of the convolution operation applied at every pixel using the convolution operation - 
            
            $$
            \begin{bmatrix}
                0 & -1 & 0\\
                -1 & 5 & -1\\
                0 & -1 & 0
            \end{bmatrix}
            $$
        
        - Blur (B) - The following kernel is used as part of the convolution operation applied at every pixel using the convolution operation - 
            
            $$
            \frac{1}{9}
            \begin{bmatrix}
                1 & 1 & 1\\
                1 & 1 & 1\\
                1 & 1 & 1
            \end{bmatrix}
            $$
        
        - Edge Detection (E) - The following kernel is used as part of the convolution operation applied at every pixel using the convolution operation - 
            
            $$
            \begin{bmatrix}
                -1 & -1 & -1\\
                -1 & 8 & -1\\
                -1 & -1 & -1
            \end{bmatrix}
            $$

        - Emboss (M) - The following kernel is used as part of the convolution operation applied at every pixel using the convolution operation -  

            $$
            \begin{bmatrix}
                -1 & -1 & 0\\
                -1 & 0 & 1\\
                0 & 1 & 1
            \end{bmatrix}
            $$


    - **Note that these identifiers for each effect are case sensitive.**
    - Convolution is perform on a 2D image grid by convolving the kernel across the image. The convolution operation can be thought of as a sliding window computation over the entire image. The computation being performed is the frobenius inner product which is the sum over element wise products between the kernel and patch of image overlapping with the kernel. **Zero padding** is used at the edges of the image. The convolution operation can be seen below - 

        $$y[m,n] = x[m,n] \ast h[m,n] = \sum_{j=-\infty}^{\infty} \sum_{i=-\infty}^{\infty} x[i,j] h[m - i, n - j]$$

4. Run Modes:
    The following modes can be used with the identifier specified in parenthesis (for sequential mode, no mode is specified when running the program).

    - sequential - If no mode is specified, then the sequential version is run.

    - work balancing (wb) - This runs the editor in the parallel work balancing mode, which balances the tasks given to each worker based on the balance threshold specified. The number of threads and the balance threshold must be specified in this mode. Each thread spawned will work on image tasks (including all of the effects for the image task). 

    - work stealing (ws) - This runs the editor in the parallel work stealing mode, in which workers steal tasks from other workers if the worker completes all of the tasks distributed to it. The number of threads must be specified in this mode. Each thread spawned will work on image tasks (including all of the effects for the image task).  

### Running the Program - 

The editor can be run in the following way - 

The first step to run any of the following modes is - 

```shell
foo@bar:~$ cd editor
```


1. Sequential run - 

    ```shell
    foo@bar:~$ go run editor.go <image directory>
    ```

2. Parallel runs - 

    1. Work Balancing - 

        ```shell
        foo@bar:~$ go run editor.go <image directory> wb <number of threads to be spawned> <balancing threshold>
        ```

    2. Work Stealing - 

        ```shell
        foo@bar:~$ go run editor.go <image directory> ws <number of threads to be spawned>
        ```

3. Multiple Input Image Directories - 

If there a images in multiple directories within `data/in`, for example, if there was the directories `small` and `big`, we can chain directories to process using `+` - 

    ```shell
    foo@bar:~$ go run editor.go small+big ws <number of threads to be spawned>
    ```


### Benchmarking the Program - 

**All testing has been carried out on the CS linux cluster i.e. the Peanut Cluster.**

Peanut Cluster specifications - 
    
1. Core architecture - Intel x86_64

2. Model name - Intel(R) Xeon(R) CPU E5-2420 0 @ 1.90GHz

3. Number of threads - 24

4. Operating system - ubuntu

5. OS version - 20.04.4 LTS



The program can be benchmarked using the following command - 

```shell
foo@bar:~$ sbatch benchmark_editor.sh
```

This must be run within the `benchmark directory`. Make sure to create the `slurm/out` directory inside `benchmark` directory and check the `.stdout` file for the outputs and the timings. 

The graph of speedups obtained can be seen below - 

1. Work balancing mode - 

    ![benchmarking_wb](./benchmark/Work-Balancing-speedup.png)

2. Work stealing mode - 

    ![benchmarking_ws](./benchmark/Work-Stealing-speedup.png)


The graphs will be created within the `benchmark` directory. The computation of the speedups along with the storing of each of the benchmarking timings and the plotting of the stored data happens by using `benchmark_graph.py` which is called from within `benchmark_editor.sh` (both reside in the `benchmark` directory).


The following observations can be made from the **work balancing** mode graph - 

1. 

The following observations can be made from the **work stealing** mode graph - 

1. 

### Questions About Implementation - 
 
1. Describe the challenges you faced while implementing the system. What aspects of the system might make it difficult to parallelize? In other words, what to you hope to learn by doing this assignment?

2. What are the hotspots (i.e., places where you can parallelize the algorithm) and bottlenecks (i.e., places where there is sequential code that cannot be parallelized) in your sequential program? Were you able to parallelize the hotspots and/or remove the bottlenecks in the parallel version?

3. What limited your speedup? Is it a lack of parallelism? (dependencies) Communication or synchronization overhead? As you try and answer these questions, we strongly prefer that you provide data and measurements to support your conclusions.

4. Compare and contrast the two parallel implementations. Are there differences in their speedups?

