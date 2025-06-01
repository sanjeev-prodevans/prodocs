#!/bin/bash
set -e

OUTDIR="bin"
APPNAME="prodocs"
SRCDIR="$(pwd)"

mkdir -p "$OUTDIR"

# List of OS/ARCH combinations
targets=(
    "windows:amd64"
    "linux:amd64"
    "darwin:amd64"
    "linux:arm64"
    "darwin:arm64"
)

for target in "${targets[@]}"; do
    IFS=":" read -r GOOS GOARCH <<< "$target"
    SUFFIX=""
    if [ "$GOOS" = "windows" ]; then
        SUFFIX=".exe"
    fi
    OUTNAME="${OUTDIR}/${GOOS}_${GOARCH}_${APPNAME}${SUFFIX}"
    echo "Building $OUTNAME ..."
    env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$OUTNAME" "$SRCDIR"
done

echo "Builds complete. Binaries are in the $OUTDIR/ directory."