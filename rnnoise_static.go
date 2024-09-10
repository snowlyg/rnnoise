//go:build !dynamic
// +build !dynamic

package rnnoise

//#cgo CFLAGS: -I${SRCDIR}/include
//#cgo CXXFLAGS: -I${SRCDIR}/include
//#cgo linux,arm64 LDFLAGS: ${SRCDIR}/lib/librnnoise-linux-arm64.a -lm//#cgo darwin,arm64 LDFLAGS: ${SRCDIR}/lib/librnnoise-darwin-arm64.a -lm
import "C"
