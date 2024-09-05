#!/bin/bash

# rm -rf src
git clone https://github.com/xiph/rnnoise.git --depth=1 src
# git clone https://github.com/TeaSpeak/rnnoise-cmake --depth=1 src
cp CMakeLists.txt ./src
cd src

# fix v.2.0  fatal error: src/_kiss_fft_guts.h: No such file or directory
sed -i 's/#include "src\/_kiss_fft_guts.h"/#include "_kiss_fft_guts.h"/' ./src/dump_features.c

export GOOS=$(go env | grep GOOS | cut -d "'" -f2)
export CC="$NDK_ROOT/$NDK_VERSION/toolchains/llvm/prebuilt/$GOOS-x86_64/bin/armv7a-linux-androideabi29-clang"
export CXX="$NDK_ROOT/$NDK_VERSION/toolchains/llvm/prebuilt/$GOOS-x86_64/bin/armv7a-linux-androideabi29-clang++"

./autogen.sh
mkdir build_android
cd build_android
cmake -DCMAKE_BUILD_TYPE=Release ..
make

cp librnnoise.a ../../lib/librnnoise-android-armv7.a
