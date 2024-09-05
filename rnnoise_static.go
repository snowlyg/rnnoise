//go:build !dynamic
// +build !dynamic

package rnnoise

//#cgo CFLAGS: -I${SRCDIR}/include
//#cgo CXXFLAGS: -I${SRCDIR}/include
//#cgo android,arm LDFLAGS: ${SRCDIR}/lib/librnnoise-android-armv7.a -lm
//#cgo linux,arm LDFLAGS: ${SRCDIR}/lib/librnnoise-linux-armv7.a -lm
//#cgo linux,arm64 LDFLAGS: ${SRCDIR}/lib/librnnoise-linux-arm64.a -lm
//#cgo linux,amd64 LDFLAGS: ${SRCDIR}/lib/librnnoise-linux-x64.a -lm
//#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/lib/librnnoise-darwin-x64.a
//#cgo darwin,arm64 LDFLAGS: ${SRCDIR}/lib/librnnoise-darwin-arm64.a
//#cgo windows,amd64 LDFLAGS: ${SRCDIR}/lib/librnnoise-windows-x64.a
import "C"
