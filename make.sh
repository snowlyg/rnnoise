#!/bin/bash

current=$(pwd)
src_path=$current/src

sudo rm -rf $src_path

git clone https://github.com/xiph/rnnoise.git -b master --depth=1 src
# git clone https://github.com/TeaSpeak/rnnoise-cmake --depth=1 src
sudo chmod -R 0777 $src_path

cd $src_path

# fix v.2.0  fatal error: src/_kiss_fft_guts.h: No such file or directory
# $src_path/src/dump_features.c
# sed -i 's/#include "src\/_kiss_fft_guts.h"/#include "_kiss_fft_guts.h"/g' $src_path/src/dump_features.c

# cat $src_path/src/dump_features.c | grep "_kiss_fft_guts.h"

./autogen.sh
./configure
make
make install

export GOOS=$(go env | grep GOOS | cut -d "'" -f2)

cp $src_path/.libs/librnnoise.a $current/lib/librnnoise-$GOOS-arm64.a
