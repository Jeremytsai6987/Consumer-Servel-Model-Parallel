import subprocess
import matplotlib.pyplot as plt
import numpy as np

import matplotlib
matplotlib.use('Agg')

test_sizes = ["xsmall", "small", "medium", "large", "xlarge"]
threads_list = [2, 4, 6, 8, 12]
run_count = 5

parallel_results = {}
serial_results = {}

for test_size in test_sizes:
    serial_results[test_size] = {}
    print(f"Running serial benchmark for {test_size}...")
    times = []
    for _ in range(run_count):
        try:
            result = subprocess.run(["go", "run", "benchmark.go", "s", test_size], capture_output=True, text=True, check=True)
            time = float(result.stdout.strip())
            times.append(time)
        except subprocess.CalledProcessError as e:
            print(f"Error: {e}")
    avg_time = np.mean(times)
    serial_results[test_size] = avg_time

for test_size in test_sizes:
    parallel_results[test_size] = {}
    for threads in threads_list:
        print(f"Running benchmark for {test_size} with {threads} threads...")
        times = []
        for _ in range(run_count):
            try:
                result = subprocess.run(["go", "run", "benchmark.go", "p", test_size, str(threads)], capture_output=True, text=True, check=True)
                time = float(result.stdout.strip())
                times.append(time)
            except subprocess.CalledProcessError as e:
                print(f"Error: {e}")
        avg_time = np.mean(times)
        parallel_results[test_size][threads] = avg_time

plt.figure(figsize=(10, 6))

for test_size in test_sizes:
    serial_time = serial_results[test_size]  
    parallel_times = parallel_results[test_size]

    speedups = {threads: serial_time / parallel_times[threads] for threads in threads_list}

    plt.plot(list(speedups.keys()), list(speedups.values()), marker='o', linestyle='-', label=f'{test_size} Test Size')

plt.xlabel('Number of Threads')
plt.ylabel('Speedup')
plt.title('Speedup Comparison for Different Test Sizes')
plt.legend()
plt.grid(True)

plt.savefig('/home/jeremyyawei/Parallel_Programming/project2/project-2-Jeremytsai6987/proj2/benchmark/speedup_comparison.png')
print("Saved combined speedup graph as 'speedup_comparison.png'")
