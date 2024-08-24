package rnnoise

/*
#cgo LDFLAGS: -lrnnoise
#include <stdlib.h>
#include "rnnoise.h"
*/
import "C"
import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"
)

const FrameSize = 480

func Run() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input noisy file> <output denoised file>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Open input and output files
	f1, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer f1.Close()

	fout, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer fout.Close()

	// Initialize RNNoise
	st := C.rnnoise_create(nil)
	defer C.rnnoise_destroy(st)

	buf := make([]int16, FrameSize)
	outBuf := make([]float32, FrameSize)

	first := true

	for {
		// Read a frame of audio samples
		err := binaryRead(f1, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read input file: %v", err)
		}

		// Convert int16 to float32
		for i := range buf {
			outBuf[i] = float32(buf[i])
		}

		// Denoise the frame
		C.rnnoise_process_frame(st, (*C.float)(unsafe.Pointer(&outBuf[0])), (*C.float)(unsafe.Pointer(&outBuf[0])))

		// Convert float32 back to int16
		for i := range outBuf {
			buf[i] = int16(outBuf[i])
		}

		// Write the denoised frame to the output file, skip the first frame
		if !first {
			err = binaryWrite(fout, buf)
			if err != nil {
				log.Fatalf("Failed to write output file: %v", err)
			}
		}

		first = false
	}

	fmt.Println("Denoising completed successfully.")
}

// binaryRead reads a frame of int16 samples from the file.
func binaryRead(f *os.File, buf []int16) error {
	return binary.Read(f, binary.LittleEndian, buf)
}

// binaryWrite writes a frame of int16 samples to the file.
func binaryWrite(f *os.File, buf []int16) error {
	return binary.Write(f, binary.LittleEndian, buf)
}
