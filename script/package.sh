#!/bin/bash

set -eu
set -o pipefail
[ "$#" = "1" ] && [ "$1" = '-v' ] && set -x

OUTPUT_DIR="bin"
PACKAGES_DIR="packages"
TEMP_DIR="temp_package"
VERSION=$(git describe --tags --always --dirty="-dev")
CHECKSUMS_FILE="$PACKAGES_DIR/checksums.txt"

make -f Makefile crossbuild

rm -rf $PACKAGES_DIR $TEMP_DIR

mkdir -p $PACKAGES_DIR $TEMP_DIR

echo "" > $CHECKSUMS_FILE

for binary in $OUTPUT_DIR/ips_*_*; do
    binary_name=$(basename $binary)

    # quick start
    if [[ $binary_name == "ips_darwin_amd64" ]]; then
        cp "$binary" "$PACKAGES_DIR/ips_macos"
        echo "$(sha256sum $PACKAGES_DIR/ips_macos | sed "s|$PACKAGES_DIR/||")" >> $CHECKSUMS_FILE
    elif [[ $binary_name == "ips_windows_amd64" ]]; then
        cp "$binary" "$PACKAGES_DIR/ips_windows.exe"
        echo "$(sha256sum $PACKAGES_DIR/ips_windows.exe | sed "s|$PACKAGES_DIR/||")" >> $CHECKSUMS_FILE
    elif [[ $binary_name == "ips_linux_amd64" ]]; then
        cp "$binary" "$PACKAGES_DIR/ips_linux"
        echo "$(sha256sum $PACKAGES_DIR/ips_linux | sed "s|$PACKAGES_DIR/||")" >> $CHECKSUMS_FILE
    fi

    cp "README.md" "README_en.md" "LICENSE" $TEMP_DIR

    package_name=""
    os_arch=$(echo $binary_name | cut -d'_' -f 2-)
    if [[ $binary_name == *"_windows_"* ]]; then
        cp "$binary" "$TEMP_DIR/ips.exe"
        package_name="ips_${VERSION}_${os_arch}.zip"
        zip -j "$PACKAGES_DIR/$package_name" -r $TEMP_DIR/*
    else
        cp "$binary" "$TEMP_DIR/ips"
        package_name="ips_${VERSION}_${os_arch}.tar.gz"
        tar -czf "$PACKAGES_DIR/$package_name" -C $TEMP_DIR .
    fi

    rm -rf $TEMP_DIR/*

    if [[ ! -z "$package_name" ]]; then
        echo "$(sha256sum $PACKAGES_DIR/$package_name | sed "s|$PACKAGES_DIR/||")" >> $CHECKSUMS_FILE
    fi

done

rm -rf $TEMP_DIR

echo "📦 All packages and their sha256 checksums have been created in $PACKAGES_DIR/"