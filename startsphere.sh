#!/bin/bash

# Check if three arguments are provided
if [ "$#" -ne 3 ]; then
    echo "Error: Please provide three arguments."
    exit 1
fi

# Start the spheres program with the provided arguments
./spheres "$2" "$3" &

# Exit the script
exit 0
