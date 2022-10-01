#!/bin/sh

# remove previous generated binaries
cd ./package
rm -rf bin/

# create binaries
yarn package

# create target directories
rm -rf tmp
rm -rf release
mkdir release
mkdir tmp

# create release zips
function createZip(){
  echo "\033[0;36mcreating zip for $1 \033[0m"
  mv bin/aws-profiles-"$1" tmp/aws-profiles
  zip -j release/aws-profiles-"$1".zip tmp/aws-profiles
  rm tmp/aws-profiles
}

createZip "linux-arm64"
createZip "linux-x64"
createZip "macos-arm64"
createZip "macos-x64"
createZip "win-x64.exe"

rm -rf bin
rm -rf tmp
