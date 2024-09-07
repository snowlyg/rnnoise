package rnnoise

// #include <stdlib.h>
// #include "rnnoise.h"
import "C"
import (
	"unsafe"
)

const MAX_FRAME_SIZE = 16384
const FrameSize = 480

func ProcessAudio(inputSamples []byte) {
	// Create a new RNNoise state
	st := C.rnnoise_create(nil)
	// Destroy the RNNoise state
	defer C.rnnoise_destroy(st)

	// inBuf := make([]byte, FrameSize)
	outBuf := make([]float32, FrameSize)
	for {
		if len(inputSamples) < FrameSize {
			break
		}

		// Convert int16 to float32
		for i := range inputSamples {
			outBuf[i] = float32(inputSamples[i])
		}
		C.rnnoise_process_frame(st, (*C.float)(unsafe.Pointer(&outBuf[0])), (*C.float)(unsafe.Pointer(&outBuf[0])))
		for i := range outBuf {
			inputSamples[i] = byte(outBuf[i])
		}
	}
}
