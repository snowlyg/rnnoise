package rnnoise

// #include <stdlib.h>
// #include "rnnoise.h"
import "C"
import (
	"unsafe"
)

func ProcessAudio(inputSamples, outputSamples []byte) {
	// Create a new RNNoise state
	st := C.rnnoise_create(nil)
	// Destroy the RNNoise state
	defer C.rnnoise_destroy(st)

	// Apply RNNoise noise reduction to each frame
	numFrames := len(inputSamples) / 4 // Assuming 16-bit PCM format (2 bytes per channel)
	for i := 0; i < numFrames; i++ {
		frameStart := i * 4
		// frameEnd := frameStart + 4

		// Convert bytes to int16 for RNNoise
		samples := *(*[2]int16)(unsafe.Pointer(&inputSamples[frameStart]))

		// Apply RNNoise

		C.rnnoise_process_frame(st, (*C.float)(unsafe.Pointer(&samples[0])), (*C.float)(unsafe.Pointer(&samples[0])))
		// Convert int16 back to bytes for output
		// copy(outputSamples[frameStart:frameEnd], *(*[4]byte)(unsafe.Pointer(&samples)))
	}

}
