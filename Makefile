# from: 
# https://github.com/bcicen/ctop/blob/master/Makefile
# https://github.com/schollz/find/blob/master/Makefile

# TODO_
# Try go dep for dependencies and see if gopherjs is stil happy. 
# Add git branching.

NAME=ctop
VERSION=$(shell cat VERSION)
BUILD_TIME=`date +%FT%T%z`
BUILD=$(shell git rev-parse --short HEAD)
LD_FLAGS="-w -X main.version=$(VERSION) -X main.build=$(BUILD) -X main.BuildTime=${BUILD_TIME}"

setup-dev-env:

	## Needs for QT


	## Needed for Mobile development
	# Install Java, Android NDK and SDK, and gomobile
	# Assumes your on OSX and using Brew. 
	brew cask list

	# Install the standard Java
	brew install Caskroom/cask/java
	
	# Install android-sdk-tools, platform-tools, and build-tools.
	brew install Caskroom/cask/android-sdk
	
	# Here are the docs...
	#	https://developer.android.com/studio/command-line/sdkmanager.html
	# 	https://developer.android.com/studio/command-line/avdmanager.html
	# You should have this available on your bash:
	# 	$ANDROID_TOOLS/bin/avdmanager
	# 	$ANDROID_TOOLS/bin/sdkmanager

	# List sdk tools
	sdkmanager --list 

	# Install android-ndk, using the android sdk tools (this is the new way wih the Android SDK.)
	# NOTE. This can take 10 minutes, and there is no stdout to tell you its happening. So just wait.
	$ANDROID_TOOLS/bin/sdkmanager "ndk-bundle"
	
	# Accept licenses (will ask you to type "y")
	$ANDROID_TOOLS/bin/sdkmanager --licenses

	# Install gomobile. https://github.com/golang/go/wiki/Mobile
	go get golang.org/x/mobile/cmd/gomobile
	gomobile init # Will take a few minutes with no feedback.

	# Test all works with a simple build
	go get -d golang.org/x/mobile/example/basic
	gomobile build -target=android golang.org/x/mobile/example/basic

	## For IOS you have to add your own ID (https://github.com/golang/mobile/blob/master/cmd/gomobile/build_iosapp.go#L34)
	code $GOPATH/src/golang.org/x/mobile/cmd/gomobile/build_iosapp.go
	## Change BundleID and then reinstall:
	go install golang.org/x/mobile/cmd/gomobile
	gomobile build -target=ios golang.org/x/mobile/example/basic
	
setup-dev-emulator:

	# ALEX: this is commneted out because its not workng yet, but close.

	# This sets up Android and IOS emulators
	# This assumes you have run the "make setup-dev-env"
	
	# Need to wrap Gradle so we dont need to screw around with Android Studio
	# https://github.com/echocat/gradle-golang-plugin
	
	java -version
	
	# This seems to have some good stuff.
	#go get github.com/bitrise-steplib/steps-start-android-emulator
	
	# Looks useful: https://github.com/bitrise-tools/go-android/...
	#go get github.com/bitrise-tools/go-android/...
	
	# List what we have now.
	#sdkmanager --list --verbose
	
	# Updater (Only run this when your sure)
	#sdkmanager --update
	
	## Install via AVD
	# $ANDROID_TOOLS/bin/avdmanager
	# $ANDROID_TOOLS/bin/avdmanager "create" "avd" "--force" "--name" "android-21" "--target" "android-21" "--abi" "armeabi-v7a"
	#android "create" "avd" "--force" "--name" "android-21" "--target" "android-21" "--abi" "armeabi-v7a"
	
	# Start emulator 
	# emulator
	

clean:
	rm -rf build/

build:
	glide install
	CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o ctop

build-dev:
	go build -ldflags "-w -X main.version=$(VERSION)-dev -X main.build=$(BUILD)"

build-all:
	mkdir -p build
	GOOS=darwin GOARCH=amd64 go build -tags release -ldflags $(LD_FLAGS) -o build/ctop-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -tags release -ldflags $(LD_FLAGS) -o build/ctop-$(VERSION)-linux-amd64
	GOOS=linux GOARCH=arm go build -tags release -ldflags $(LD_FLAGS) -o build/ctop-$(VERSION)-linux-arm


.PHONY: build