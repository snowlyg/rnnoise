package rnnoise

// #include <stdlib.h>
// #include "rnnoise.h"
import "C"
import (
	"bytes"
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

// Process
func Process(sampleCount []byte) []byte {

	if len(sampleCount) < FrameSize {
		println("input < 480")
		return sampleCount
	}

	// Create a new RNNoise state
	st := C.rnnoise_create(nil)
	// Destroy the RNNoise state
	defer C.rnnoise_destroy(st)

	piBuffer := bytes.NewReader(sampleCount)

	inputTmp := make([]int16, FrameSize)
	outTmp := make([]float32, FrameSize)

	binaryRead(piBuffer, inputTmp)

	for i := 0; i < FrameSize; i++ {
		outTmp[i] = float32(inputTmp[i])
	}

	C.rnnoise_process_frame(st, (*C.float)(unsafe.Pointer(&outTmp[0])), (*C.float)(unsafe.Pointer(&outTmp[0])))

	for i := 0; i < FrameSize; i++ {
		inputTmp[i] = int16(outTmp[i])
	}

	buf := new(bytes.Buffer)
	binaryWrite(buf, inputTmp)

	out := make([]byte, len(sampleCount))
	m, err := buf.Read(out)
	if err == io.EOF {
		println("EOF:", err.Error())
		// break
	}

	return out[:m]

}

// ProcessFile
func ProcessFile(inputFile string) {
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

	sampleCount := make([]byte, FrameSize*2)

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

			// piBuffer := bytes.NewReader(sampleCount[:x])

			// inputTmp := make([]int16, FrameSize)
			// outTmp := make([]float32, FrameSize)

			// binaryRead(piBuffer, inputTmp)

			// for i := 0; i < FrameSize; i++ {
			// 	outTmp[i] = float32(inputTmp[i])
			// }

			// if len(inputTmp) < FrameSize {
			// 	println("input < 480")
			// 	break
			// }

			// C.rnnoise_process_frame(st, (*C.float)(unsafe.Pointer(&outTmp[0])), (*C.float)(unsafe.Pointer(&outTmp[0])))

			// for i := 0; i < FrameSize; i++ {
			// 	inputTmp[i] = int16(outTmp[i])
			// }

			// buf := new(bytes.Buffer)
			// binaryWrite(buf, inputTmp)

			// out := make([]byte, x)
			// m, err := buf.Read(out)
			// if err == io.EOF {
			// 	println("EOF:", err.Error())
			// 	break
			// }
			// println("x", x)
			// println("m", m)

			out := Process(sampleCount[:x])

			w.Write(out)
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
func binaryRead(f io.Reader, buf []int16) error {
	return binary.Read(f, binary.LittleEndian, buf)
}

// binaryWrite writes a frame of int16 samples to the file.
func binaryWrite(f io.Writer, buf []int16) error {
	return binary.Write(f, binary.LittleEndian, buf)
}
