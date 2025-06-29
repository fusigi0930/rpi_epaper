export CROSS_COMPILE=~/toolchain/aarch64/bin/aarch64-none-linux-gnu-
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=arm64
export CC=${CROSS_COMPILE}gcc
export SYSROOT=/home/chuyuan/sysroot/rpi_b2qt

cmake -B build_rpi -DCMAKE_BUILD_TYPE=release -DTARGET_TAG=rpi
cmake --build build_rpi -j8 -v
