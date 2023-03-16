#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e
# Print commands and their arguments as they are executed.
set -x

appName="JDC-Monitor"

make -f Makefile build
cd build
mkdir -p compress

for i in $(find . -type f -name "$appName-linux-*"); do
  # why? Because some target platforms seem to have issues with upx compression
  if [[ "$i" == "./$appName-linux-mips64"* || "$i" == "./$appName-linux-riscv64"* ]]; then
    cp "$i" $appName
  else
    upx --lzma --best "$i" -o $appName
  fi
  tar -czvf compress/"$i".tar.gz $appName
  rm -f $appName
done

for i in $(find . -type f -name "$appName-darwin-*"); do
  if [[ "$i" == "./$appName-darwin-arm64"* ]]; then
    cp "$i" $appName
  else
    upx --lzma --best "$i" -o $appName
  fi
  tar -czvf compress/"$i".tar.gz $appName
  rm -f $appName
done

for i in $(find . -type f -name "$appName-freebsd-*"); do
  cp "$i" $appName
  tar -czvf compress/"$i".tar.gz $appName
  rm -f $appName
done

for i in $(find . -type f -name "$appName-windows-*"); do
  if [[ "$i" == "./$appName-windows-arm64"* ]]; then
    cp "$i" $appName.exe
  else
    upx --lzma --best "$i" -o $appName.exe
  fi
  zip compress/$(echo $i | sed 's/\.[^.]*$//').zip $appName.exe
  rm -f $appName.exe
done

cd compress
find . -type f -print0 | xargs -0 md5sum >md5.txt
cat md5.txt
cd ../..
