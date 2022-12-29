import subprocess
import sys
import matplotlib

import matplotlib.pyplot as plt

test_sizes = ["small", "mixture", "big"]
threadNum = ["2", "4", "6", "8", "12"]
matplotlib.use('Agg')

# collect sequential data
sequential_time = {}
for sz in test_sizes:
    print("start running version: s, " + "size: " + sz)
    total = 0
    for i in range(5):
        p = subprocess.Popen(['go', 'run', '../editor/editor.go', sz],
                             stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        stdout, stderr = p.communicate()
        time = float(stdout.decode("utf-8").strip())
        print("sequential version, size: %s, number of run %d"%(sz, i))
        total = total + time
    sequential_time[sz] = total / 5
print(sequential_time)

# collect parallel data for both strategies
for strategy in ["pipeline", "bsp"]:
    res = {}
    for sz in test_sizes:
        data = []
        print("start running version:" + strategy + ", size: " + sz)
        for tn in threadNum:
            print("thread Num: " + tn)
            total = 0
            for i in range(5):
                p = subprocess.Popen(['go', 'run', '../editor/editor.go', sz, strategy, tn],
                                    stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
                stdout, stderr = p.communicate()
                time = float(stdout.decode("utf-8").strip())
                print("%s: %s, threadNum: %d, runNumber: %d"%(strategy, sz, int(tn), i))
                total = total + time
            average = total / 5
            print(average)
            speedup = sequential_time[sz] / average
            data.append(speedup)
        res[sz] = data
    print(res)

    # Plot the graph
    plt.figure()
    for i, size in enumerate(test_sizes):
        plt.plot(threadNum,
            res[size],
            label=test_sizes[i],
            linestyle='dashed',
            linewidth=3,
            marker='o',
            markersize=9)

    plt.xlabel('Number of threads')
    plt.ylabel('Speed Up')
    plt.title('Speedup Graph')
    plt.legend()
    plt.savefig(f"speedup-{strategy}.png")

