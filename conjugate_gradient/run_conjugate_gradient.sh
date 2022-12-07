#!/bin/bash

# Get n from the user
echo "Enter n:"
read n

# Run the conjugate gradient go script
go run simulator/simulator.go $n

# Run the visualization python script
python3 visualize.py