#!/bin/bash
# Build releases for all platforms

set -e

VERSION=${1:-dev}
OUTPUT_DIR="dist"

echo "Building Twitter CLI $VERSION for all platforms..."
echo ""

# Clean and create output directory
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Platforms to build for
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    IFS='/' read -r -a PARTS <<< "$PLATFORM"
    GOOS="${PARTS[0]}"
    GOARCH="${PARTS[1]}"
    
    OUTPUT_NAME="twt-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags="-X 'main.Version=$VERSION' -X 'main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
        -o "${OUTPUT_DIR}/${OUTPUT_NAME}" \
        .
    
    echo "✓ Built ${OUTPUT_NAME}"
done

echo ""
echo "Build complete! Binaries are in $OUTPUT_DIR/"
echo ""

# Create checksums
cd "$OUTPUT_DIR"
sha256sum twt-* > SHA256SUMS
cd ..

echo "Checksums saved to $OUTPUT_DIR/SHA256SUMS"