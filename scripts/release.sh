#!/bin/bash

VERSION="0.1.0"
BINARY_NAME="ginvader"
PLATFORMS="darwin/amd64"

# Create release directory
rm -rf release
mkdir -p release

# Build for each platform
for PLATFORM in $PLATFORMS; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  
  OUTPUT_NAME=$BINARY_NAME
  if [ $GOOS = "windows" ]; then
    OUTPUT_NAME+='.exe'
  fi

  echo "Building for $GOOS/$GOARCH..."
  GOOS=$GOOS GOARCH=$GOARCH go build -o release/$OUTPUT_NAME ./cmd/ginvader
  
  # Create tar.gz
  cd release
  tar -czf ${BINARY_NAME}_${VERSION}_${GOOS}_${GOARCH}.tar.gz $OUTPUT_NAME
  rm $OUTPUT_NAME
  cd ..
done

# Generate SHA256 for the archive
cd release
shasum -a 256 ${BINARY_NAME}_${VERSION}_darwin_amd64.tar.gz > SHA256SUMS
cd ..
