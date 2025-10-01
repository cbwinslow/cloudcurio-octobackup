#!/bin/bash

# Run the CloudCurio Code Review Crew

# Check if OPENAI_API_KEY is set
if [ -z "$OPENAI_API_KEY" ]; then
    echo "Error: OPENAI_API_KEY environment variable is not set."
    echo "Please set your OpenAI API key and try again."
    exit 1
fi

# Navigate to the crew directory and run the main script
cd "$(dirname "$0")/crew"
echo "Starting CloudCurio Code Review Crew..."
python3 main.py