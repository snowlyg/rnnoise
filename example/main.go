package main

// #include <stdlib.h>
// #include "rnnoise.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"unsafe"

	"github.com/gen2brain/malgo"
)

const FrameSize = 480

func b2f32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func f322b(f float32) []byte {
	bits := math.Float32bits(f)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input noisy file> <output denoised file>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Open input and output files
	f1, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer f1.Close()

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	sampleCount := make([]byte, FrameSize*4)

	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()
	// Create a new RNNoise state
	st := C.rnnoise_create(nil)
	// Destroy the RNNoise state
	defer C.rnnoise_destroy(st)
	go func() {
		for {
			x, err := f1.Read(sampleCount)
			if err == io.EOF {
				println("EOF:", err.Error())
				break
			}

			inputTmp := make([]float32, FrameSize)

			inputIndex := 0
			for i := 0; i < FrameSize; i++ {
				inputTmp[i] = b2f32(sampleCount[inputIndex : inputIndex+4])
				inputIndex += 4
			}

			out := make([]byte, 0)
			if len(inputTmp) < FrameSize {
				println("input < 480")
				return
			}

			// Convert int16 to float32
			C.rnnoise_process_frame(st, (*C.float)(unsafe.Pointer(&inputTmp[0])), (*C.float)(unsafe.Pointer(&inputTmp[0])))

			outputIndex := 0
			for i := 0; i < FrameSize; i++ {
				out = append(out, f322b(inputTmp[i])...)
				outputIndex += 4
			}
			println("out", len(out))
			w.Write(out[:x])
		}
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = 48000
	deviceConfig.Alsa.NoMMap = 1

	// This is the function that's used for sending more data to the device for playback.
	onSamples := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		io.ReadFull(r, pOutputSample)
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onSamples,
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer device.Uninit()

	err = device.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Press Enter to quit...")
	fmt.Scanln()
}

// binaryRead reads a frame of int16 samples from the file.
func binaryRead(f *os.File, buf []int16) error {
	return binary.Read(f, binary.LittleEndian, buf)
}

// binaryWrite writes a frame of int16 samples to the file.
func binaryWrite(f *os.File, buf []int16) error {
	return binary.Write(f, binary.LittleEndian, buf)
}
