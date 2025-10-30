#!/bin/bash

# Path to the binary
BINARY_PATH=


# Find running PIDs of botsupervisor
PIDS=$(pgrep -x "botsupervisor")

if [ -n "$PIDS" ]; then
    echo "$(date): Terminating running botsupervisor process(es): $PIDS"
    # Kill the running process(es)
    kill $PIDS

    # Wait for processes to terminate gracefully (up to 5 seconds)
    for i in {1..5}; do
        if pgrep -x "botsupervisor" > /dev/null; then
            sleep 1
        else
            break
        fi
    done

    # Force kill if still running
    if pgrep -x "botsupervisor" > /dev/null; then
        echo "$(date): Force killing botsupervisor process(es): $PIDS"
        kill -9 $PIDS
    fi
else
    echo "$(date): No running botsupervisor processes found."
fi

echo "$(date): Starting botsupervisor."
# Start the binary in background
"$BINARY_PATH" &

