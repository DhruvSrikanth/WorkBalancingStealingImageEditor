import subprocess
import matplotlib.pyplot as plt

def get_sequential_times(problem_sizes, average_over, run_command):
    '''
    Get the sequential runtimes for the benchmark graphs
    Args:
        problem_sizes: The size of the problem.
        average_over: The number of times to run the benchmark and average timings over.
        run_command: The command to run the benchmark with.
    Returns:
        A dictionary of the sequential runtimes for each benchmark graph.
    '''
    encoding = 'utf-8'
    sequential_times = {problem_size: 0 for problem_size in problem_sizes}
    print('Running sequential benchmarks...')
    print('--------------------------------')
    for problem_size in problem_sizes:
        print('--------------------------------')
        for i in range(average_over):
            command = run_command.format(problem_size)
            time = subprocess.check_output(command, shell=True)
            time = str(time, encoding).strip().replace('\n', '')
            sequential_times[problem_size] += float(time)
            print(f'Using command - {command}. Mode: sequential, Problem size: {problem_size}, Iteration: {i + 1}, Time: {time}s')
        sequential_times[problem_size] /= average_over
        print('--------------------------------')
    return sequential_times

def get_parallel_times(problem_sizes, thread_nums, average_over, run_command, mode):
    '''
    Get the parallel runtimes for the benchmark graphs
    Args:
        problem_sizes: The size of the problem.
        thread_nums: The number of threads to run the benchmark with.
        average_over: The number of times to run the benchmark and average timings over.
        run_command: The command to run the benchmark with.
        mode: The mode of the benchmark.
    Returns:
        A dictionary of the parallel runtimes for each benchmark graph.
    '''
    encoding = 'utf-8'
    parallel_times = {problem_size: {str(thread_num): 0 for thread_num in thread_nums} for problem_size in problem_sizes}
    print('Running parallel benchmarks...')
    print('--------------------------------')
    for problem_size in problem_sizes:
        print('--------------------------------')
        for thread_num in thread_nums:
            print('--------------------------------')
            for i in range(average_over):
                command = run_command.format(problem_size, thread_num)
                time = subprocess.check_output(command, shell=True)
                time = str(time, encoding).strip().replace('\n', '')
                parallel_times[problem_size][str(thread_num)] += float(time)
                print(f'Using command - {command}. Mode: {mode}, Problem size: {problem_size}, Thread num: {thread_num}, Iteration: {i + 1}, Time: {time}s')
            parallel_times[problem_size][str(thread_num)] /= average_over
            print('--------------------------------')
        print('--------------------------------')
    return parallel_times

def get_speedup(sequential_time, parallel_time):
    '''
    Get the speedup of a parallel run compared to a sequential run.
    Args:
        sequential_time: The time taken for the sequential run.
        parallel_time: The time taken for the parallel run.
    Returns:
        The speedup of the parallel run.
    '''
    return sequential_time / parallel_time

def get_speedups(sequential_times, parallel_times):
    '''
    Get the speedups of the parallel runs compared to the sequential runs.
    Args:
        sequential_times: The times taken for the sequential runs.
        parallel_times: The times taken for the parallel runs.
    Returns:
        A dictionary of the speedups for each benchmark graph.
    '''
    speedups = {problem_size: {str(thread_num): 0 for thread_num in parallel_times[problem_size]} for problem_size in parallel_times}
    for problem_size in parallel_times:
        for thread_num in parallel_times[problem_size]:
            speedups[problem_size][thread_num] = get_speedup(sequential_times[problem_size], parallel_times[problem_size][thread_num])
    return speedups

def plot_speedups(num_threads, speedups, mode):
    '''
    Plot the speedups for each benchmark.
    Args:
        num_threads: The number of threads used for the parallel runs.
        speedups: The speedups for each benchmark graph.
        mode: The mode of the benchmark.
    '''
    for problem_size in speedups:
        plt.plot(num_threads, list(speedups[problem_size].values()), label=f'{problem_size} data points')

    plt.xlabel('Number of threads')
    plt.ylabel('Speedup')
    plt.title(f'Sparse Matrix Solver using Conjugate Gradient ({mode})')
    plt.legend(loc='best')
    plt.tight_layout()
    plt.grid()
    plt.savefig(f'{mode}-speedup.png')
    plt.clf()

if __name__ == '__main__':
    # The different problem sizes to benchmark.
    problem_sizes = [
        100,
        500,
        1000,
    ]

    # The number of threads to benchmark with.
    thread_nums = [
        2, 4, 6, 8, 12
    ]
    
    # The number of times to run the benchmark and average timings over.
    average_over = 5

    # The command to run the benchmark with.
    run_command_sequential = 'go run ../simulator/simulator.go {}'

    # Get the sequential runtimes.
    sequential_times = get_sequential_times(problem_sizes=problem_sizes, average_over=average_over, run_command=run_command_sequential)
    print(f'Sequential times: {sequential_times}\n')
    
    run_command_parallel = 'go run ../simulator/simulator.go {} {}'

    mode = 'MapReduce'
    # Get the parallel runtimes.
    parallel_times = get_parallel_times(problem_sizes=problem_sizes, thread_nums=thread_nums, average_over=average_over, run_command=run_command_parallel, mode=mode)
    print(f'Parallel times for {mode}: {parallel_times}\n')

    # Get the speedups.
    speedups = get_speedups(sequential_times=sequential_times, parallel_times=parallel_times)
    print(f'Speedups for {mode}: {speedups}\n')

    # Plot the speedups.
    print(f'Plotting speedups for {mode}...')
    plot_speedups(num_threads=thread_nums, speedups=speedups, mode=mode)
    print('Done!')