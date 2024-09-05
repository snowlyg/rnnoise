git clone https://github.com/xiph/rnnoise.git --depth=1 src

go build -o denoise main.go
./denoise input_noisy.wav output_denoised.wav
