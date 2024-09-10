# rnnoise

## How to use

```
git clone https://github.com/snowlyg/rnnoise.git --depth=1

cd rnnoise/example
go build -o denoise main.go 

./main 

./denoise origin.wav denoised.wav
```

## build lib

```
git clone https://github.com/xiph/rnnoise.git --depth=1 src

chmod -R 0777 ./src

cd ./src

./autogen.sh
./configure  
make
make install

```

## use for android

[RnnoiseAndroidPlugin](https://github.com/snowlyg/RnnoiseAndroidPlugin.git)