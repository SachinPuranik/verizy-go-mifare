#!/bin/bash
TARGET_USER=pi
TARGET_HOST=192.168.0.196
TARGET_DIR=cardtest
TARGET_BINARY=main
ARM_VERSION=6
 
# Executable name is assumed to be same as current directory name
#EXECUTABLE=${PWD##*/} 
EXECUTABLE=main
 
echo "Building for Raspberry Pi..."
env GOOS=linux GOARCH=arm GOARM=$ARM_VERSION go bui ld
 
echo "Uploading to Raspberry Pi..."
#scp -i ~/.ssh/dev-key $EXECUTABLE $TARGET_USER@$TARGET_HOST:$TARGET_DIR/$EXECUTABLE
scp  $EXECUTABLE $TARGET_USER@$TARGET_HOST:~/$TARGET_DIR/$EXECUTABLE