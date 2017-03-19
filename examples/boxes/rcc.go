package main

/*
#cgo +build windows,386 windows,amd64 LDFLAGS: -L/usr/lib/mxe/usr/x86_64-w64-mingw32.shared/qt5/lib -lQt5Core

#cgo +build !ios,darwin,amd64 LDFLAGS: -F/home/influx6/Labs/qt5/5.8/clang_64/lib -framework QtCore

#cgo +build linux,amd64 LDFLAGS: -Wl,-rpath,/home/influx6/Labs/qt5/5.8/gcc_64/lib -L/home/influx6/Labs/qt5/5.8/gcc_64/lib -lQt5Core


#cgo +build android,linux,arm LDFLAGS: -L/home/influx6/Labs/qt5/5.8/android_armv7/lib -lQt5Core


#cgo +build ios,darwin,amd64 LDFLAGS: -headerpad_max_install_names -stdlib=libc++ -Wl,-syslibroot,/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/SDKs/ -mios-simulator-version-min=7.0 -arch x86_64
#cgo +build ios,darwin,amd64 LDFLAGS: -L/home/influx6/Labs/qt5/5.8/ios/plugins/platforms -lqios -framework Foundation -framework UIKit -framework QuartzCore -framework AudioToolbox -framework AssetsLibrary -L/home/influx6/Labs/qt5/5.8/ios/lib -framework MobileCoreServices -framework CoreFoundation -framework OpenGLES -framework CoreText -framework CoreGraphics -framework Security -framework SystemConfiguration -framework CoreBluetooth -lQt5FontDatabaseSupport -lQt5GraphicsSupport -lQt5ClipboardSupport -lqtfreetype -L/home/influx6/Labs/qt5/5.8/ios/plugins/imageformats -lqgif -lqicns -lqico -lqjpeg -lqmacjp2 -framework ImageIO -lqtga -lqtiff -lqwbmp -lqwebp -lqtlibpng -lqtharfbuzz -lm -lz -lqtpcre -lQt5Core -lQt5Widgets -lQt5Gui

#cgo +build ios,darwin,386 LDFLAGS: -headerpad_max_install_names -stdlib=libc++ -Wl,-syslibroot,/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/SDKs/ -mios-simulator-version-min=7.0 -arch i386
#cgo +build ios,darwin,386 LDFLAGS: -L/home/influx6/Labs/qt5/5.8/ios/plugins/platforms -lqios -framework Foundation -framework UIKit -framework QuartzCore -framework AudioToolbox -framework AssetsLibrary -L/home/influx6/Labs/qt5/5.8/ios/lib -framework MobileCoreServices -framework CoreFoundation -framework OpenGLES -framework CoreText -framework CoreGraphics -framework Security -framework SystemConfiguration -framework CoreBluetooth -lQt5FontDatabaseSupport -lQt5GraphicsSupport -lQt5ClipboardSupport -lqtfreetype -L/home/influx6/Labs/qt5/5.8/ios/plugins/imageformats -lqgif -lqicns -lqico -lqjpeg -lqmacjp2 -framework ImageIO -lqtga -lqtiff -lqwbmp -lqwebp -lqtlibpng -lqtharfbuzz -lm -lz -lqtpcre -lQt5Core -lQt5Widgets -lQt5Gui

#cgo +build ios,darwin,arm64 LDFLAGS: -headerpad_max_install_names -stdlib=libc++ -Wl,-syslibroot,/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/ -miphoneos-version-min=7.0 -arch arm64
#cgo +build ios,darwin,arm64 LDFLAGS: -L/home/influx6/Labs/qt5/5.8/ios/plugins/platforms -lqios -framework Foundation -framework UIKit -framework QuartzCore -framework AudioToolbox -framework AssetsLibrary -L/home/influx6/Labs/qt5/5.8/ios/lib -framework MobileCoreServices -framework CoreFoundation -framework OpenGLES -framework CoreText -framework CoreGraphics -framework Security -framework SystemConfiguration -framework CoreBluetooth -lQt5FontDatabaseSupport -lQt5GraphicsSupport -lQt5ClipboardSupport -lqtfreetype -L/home/influx6/Labs/qt5/5.8/ios/plugins/imageformats -lqgif -lqicns -lqico -lqjpeg -lqmacjp2 -framework ImageIO -lqtga -lqtiff -lqwbmp -lqwebp -lqtlibpng -lqtharfbuzz -lm -lz -lqtpcre -lQt5Core -lQt5Widgets -lQt5Gui

#cgo +build ios,darwin,arm LDFLAGS: -headerpad_max_install_names -stdlib=libc++ -Wl,-syslibroot,/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/ -miphoneos-version-min=7.0 -arch armv7
#cgo +build ios,darwin,arm LDFLAGS: -L/home/influx6/Labs/qt5/5.8/ios/plugins/platforms -lqios -framework Foundation -framework UIKit -framework QuartzCore -framework AudioToolbox -framework AssetsLibrary -L/home/influx6/Labs/qt5/5.8/ios/lib -framework MobileCoreServices -framework CoreFoundation -framework OpenGLES -framework CoreText -framework CoreGraphics -framework Security -framework SystemConfiguration -framework CoreBluetooth -lQt5FontDatabaseSupport -lQt5GraphicsSupport -lQt5ClipboardSupport -lqtfreetype -L/home/influx6/Labs/qt5/5.8/ios/plugins/imageformats -lqgif -lqicns -lqico -lqjpeg -lqmacjp2 -framework ImageIO -lqtga -lqtiff -lqwbmp -lqwebp -lqtlibpng -lqtharfbuzz -lm -lz -lqtpcre -lQt5Core -lQt5Widgets -lQt5Gui


#cgo +build sailfish_emulator,linux,386 LDFLAGS: -Wl,-rpath,/usr/share/harbour-boxes/lib -Wl,-rpath-link,/srv/mer/targets/SailfishOS-i486/usr/lib -Wl,-rpath-link,/srv/mer/targets/SailfishOS-i486/lib -L/srv/mer/targets/SailfishOS-i486/usr/lib -L/srv/mer/targets/SailfishOS-i486/lib -lQt5Core
#cgo +build sailfish,linux,arm LDFLAGS: -Wl,-rpath,/usr/share/harbour-boxes/lib -Wl,-rpath-link,/srv/mer/targets/SailfishOS-armv7hl/usr/lib -Wl,-rpath-link,/srv/mer/targets/SailfishOS-armv7hl/lib -L/srv/mer/targets/SailfishOS-armv7hl/usr/lib -L/srv/mer/targets/SailfishOS-armv7hl/lib -lQt5Core

#cgo +build asteroid,linux,arm LDFLAGS: -Wl,-rpath-link,/usr/lib -Wl,-rpath-link,/lib -L/usr/lib -L/lib -lQt5Core


#cgo +build rpi1,linux,arm LDFLAGS: -Wl,-rpath-link,/home/influx6/raspi/sysroot/opt/vc/lib -Wl,-rpath-link,/home/influx6/raspi/sysroot/usr/lib/arm-linux-gnueabihf -Wl,-rpath-link,/home/influx6/raspi/sysroot/lib/arm-linux-gnueabihf -Wl,-rpath-link,/home/influx6/Labs/qt5/5.8/rpi1/lib -mfloat-abi=hard --sysroot=/home/influx6/raspi/sysroot -Wl,-O1 -Wl,--enable-new-dtags -Wl,-z,origin -L/home/influx6/raspi/sysroot/opt/vc/lib -L/home/influx6/Labs/qt5/5.8/rpi1/lib -lQt5Core
#cgo +build rpi2,linux,arm LDFLAGS: -Wl,-rpath-link,/home/influx6/raspi/sysroot/opt/vc/lib -Wl,-rpath-link,/home/influx6/raspi/sysroot/usr/lib/arm-linux-gnueabihf -Wl,-rpath-link,/home/influx6/raspi/sysroot/lib/arm-linux-gnueabihf -Wl,-rpath-link,/home/influx6/Labs/qt5/5.8/rpi2/lib -mfloat-abi=hard --sysroot=/home/influx6/raspi/sysroot -Wl,-O1 -Wl,--enable-new-dtags -Wl,-z,origin -L/home/influx6/raspi/sysroot/opt/vc/lib -L/home/influx6/Labs/qt5/5.8/rpi2/lib -lQt5Core
#cgo +build rpi3,linux,arm LDFLAGS: -Wl,-rpath-link,/home/influx6/raspi/sysroot/opt/vc/lib -Wl,-rpath-link,/home/influx6/raspi/sysroot/usr/lib/arm-linux-gnueabihf -Wl,-rpath-link,/home/influx6/raspi/sysroot/lib/arm-linux-gnueabihf -Wl,-rpath-link,/home/influx6/Labs/qt5/5.8/rpi3/lib -mfloat-abi=hard --sysroot=/home/influx6/raspi/sysroot -Wl,-O1 -Wl,--enable-new-dtags -Wl,-z,origin -L/home/influx6/raspi/sysroot/opt/vc/lib -L/home/influx6/Labs/qt5/5.8/rpi3/lib -lQt5Core
*/
import "C"