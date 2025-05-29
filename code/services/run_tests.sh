#!/bin/bash

services=(
    "vote-service"
    "thread-service"
    "search-service"
    "popular-service"
    "community-service"
    "comment-service"
)

for service in "${services[@]}"; do
    echo "Running tests for $service..."
    cd "$service/test" || {
        echo "Failed to change directory to $service/test"
        continue
    }
    
    # Run the tests
    go test -v ./...
    
    # Return to the services directory
    cd ../..
    
    echo "----------------------------------------"
done

echo "All tests completed!" 