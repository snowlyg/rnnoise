#!/bin/bash

sudo rm -rf ./src
sudo git clone https://github.com/xiph/rnnoise.git --depth=1 src
# git clone https://github.com/TeaSpeak/rnnoise-cmake --depth=1 src
sudo cp CMakeLists.txt ./src
cd ./src

src_path=$(pwd)

# fix v.2.0  fatal error: src/_kiss_fft_guts.h: No such file or directory
sudo sed -e 's/#include "src\/_kiss_fft_guts.h"/#include "_kiss_fft_guts.h"/g' $src_path/src/dump_features.c >$src_path/src/dump_features.c

export GOOS=$(go env | grep GOOS | cut -d "'" -f2)
export CC="$NDK_ROOT/$NDK_VERSION/toolchains/llvm/prebuilt/$GOOS-x86_64/bin/armv7a-linux-androideabi29-clang"
export CXX="$NDK_ROOT/$NDK_VERSION/toolchains/llvm/prebuilt/$GOOS-x86_64/bin/armv7a-linux-androideabi29-clang++"

sudo ./autogen.sh
sudo rm -rf build_android
sudo mkdir build_android
cd build_android
sudo cmake -DCMAKE_BUILD_TYPE=Release ..
sudo make

sudo cp librnnoise.a ../../lib/librnnoise-drawin-armv7.a
