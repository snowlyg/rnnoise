//go:build !dynamic
// +build !dynamic

package rnnoise

//#cgo CFLAGS: -I${SRCDIR}/include
//#cgo CXXFLAGS: -I${SRCDIR}/include
//#cgo android,arm LDFLAGS: ${SRCDIR}/lib/librnnoise-android-armv7.a -lm
//#cgo linux,x64 LDFLAGS: ${SRCDIR}/lib/librnnoise-linux-x64.a -lm
import "C"
