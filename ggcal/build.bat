set ORI_PATH=%PATH%
set MINGW_ROOT=X:\toolchain\msys64\ucrt64\bin
set SDL2_ROOT=X:\DevelopLib\sdl2\mingw\SDL2-2.32.8\x86_64-w64-mingw32
set PATH=%PATH%;%MINGW_ROOT%

cmake -B build_win -G Ninja ^
	-DCMAKE_MAKE_PROGRAM=ninja.exe ^
	-DCMAKE_BUILD_TYPE=RELEASE

cmake --build build_win -j8

set PATH=%ORI_PATH%