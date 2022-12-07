#!/bin/bash

# Ask user which mode to run in
echo "Choose mode (Enter 1 or 2):"
echo "1. Run serial implementation"
echo "2. Run parallel implementation (MapReduce)"
read mode

# Check for invalid input
if [ $mode -ne 1 ] && [ $mode -ne 2 ]; then
    echo "Invalid input"
    exit 1
fi

# Get n from the user
echo "Enter n:"
read n

# Run the conjugate gradient go script
if [ $mode -eq 1 ]; then
    # Run the serial implementation
    go run simulator/simulator.go $n
else
    # Get the number of threads
    echo "Enter number of threads:"
    read threads

    # Run the parallel implementation
    go run simulator/simulator.go $n $threads
fi

# Run the visualization python script
python3 visualize.py

echo "Visualization inside conjugate_gradient/output/"