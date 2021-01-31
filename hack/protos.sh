#!/bin/bash

# Clean up old shite
find ./ -maxdepth 4 -type f -name "*.pb.go" -exec rm -rf {} \;

# Generate new files
files=$(find ./ -maxdepth 4 -type f -name "*.proto")
for f in $files; do
    echo -e "Compiling: $f";
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $f
done