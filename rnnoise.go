package rnnoise

/*
#cgo LDFLAGS: -lrnnoise
#include <stdlib.h>
#include "rnnoise.h"
*/
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
	"github.com/youpy/go-wav"
)

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

var st *C.DenoiseState

const FrameSize = 480

// DenoiseState
type DenoiseState struct {
	ds *C.DenoiseState
}

// NewDenoiseState
func NewDenoiseState() *DenoiseState {
	return &DenoiseState{
		// ds: C.rnnoise_create(C.rnnoise_model_from_filename(C.CString("weights_blob.bin"))),
		ds: C.rnnoise_create(nil),
	}

}

// DestoryDenoiseState
func (d *DenoiseState) DestoryDenoiseState() {
	if d.ds != nil {
		C.rnnoise_destroy(d.ds)
	}
}

// Denoise
func (d *DenoiseState) Denoise(samples []byte) []byte {
	fin := bytes.NewReader(samples)

	buf := make([]int16, FrameSize)
	binary.Read(fin, binary.LittleEndian, buf)

	buf = d.DenoiseInt16(buf)

	fout := new(bytes.Buffer)
	binary.Write(fout, binary.LittleEndian, buf)

	out := make([]byte, len(samples))
	m, _ := fout.Read(out)

	denoise := make([]byte, m)
	copy(denoise, out[:m])

	return denoise
}

// Process
func (d *DenoiseState) DenoiseInt16(inputTmp []int16) []int16 {
	tmp := make([]float32, FrameSize)
	for i := 0; i < FrameSize; i++ {
		tmp[i] = float32(inputTmp[i])
	}

	C.rnnoise_process_frame(d.ds, (*C.float)(unsafe.Pointer(&tmp[0])), (*C.float)(unsafe.Pointer(&tmp[0])))

	if len(tmp) < FrameSize {
		log.Printf("rnnoise_process_frame return len is %d < %d \n\t", len(tmp), FrameSize)
	}

	for i := 0; i < FrameSize; i++ {
		inputTmp[i] = int16(tmp[i])
	}

	return inputTmp
}

// PlayFile
func PlayFile(inputFile string) {
	// Open input and output files
	f1, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer f1.Close()
	reader := wav.NewReader(f1)
	format, err := reader.Format()
	if err != nil {
		log.Fatalf("Failed to get file Format: %v", err.Error())
	}

	log.Printf("%+v\t\n", format)

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	sampleSize := make([]byte, FrameSize*2)

	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()

	ds := NewDenoiseState()
	defer ds.DestoryDenoiseState()

	go func() {
		for {
			x, err := f1.Read(sampleSize)
			if err == io.EOF {
				println("EOF:", err.Error())
				break
			}
			// data := sampleSize[:x]
			data := ds.Denoise(sampleSize[:x])
			w.Write(data)
		}
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = uint32(format.NumChannels)
	deviceConfig.SampleRate = format.SampleRate
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
