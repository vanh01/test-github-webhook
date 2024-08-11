#!/bin/bash

# Find the process ID (PID) of the chatbot process
chatbot_pid=$(pgrep -f "githubwebhook")

# Check if the process is running
if [ -n "$chatbot_pid" ]; then
    echo "Github webhook process found with PID: $chatbot_pid"

    # Kill the chatbot process
    kill "$chatbot_pid"
    echo "Github webhook process killed."
else
    echo "Github webhook process not found."
fi

go build -o ./githubwebhook .
nohup ./githubwebhook &