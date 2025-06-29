cmake_minimum_required(VERSION 3.16)
if(LINUX)
	list(APPEND CMAKE_MODULE_PATH "${CMAKE_SOURCE_DIR}/cmake")
	project(ggcal LANGUAGES C)
endif()

if(WIN32)
	set(CMAKE_DISABLE_FIND_PACKAGE_MSVC TRUE)
	set(CMAKE_C_COMPILER_WORKS TRUE)
	set(CMAKE_CXX_COMPILER_WORKS TRUE)
	project(ggcal NONE)
endif()

#find_package(Git REQUIRED)
#execute_process(
#	COMMAND ${GIT_EXECUTABLE} describe
#	OUTPUT_VARIABLE GIT_VERSION
#	WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
#	OUTPUT_STRIP_TRAILING_WHITESPACE
#	ERROR_QUIET
#)
#
#if(NOT GIT_VERSION)
#	set(GIT_VERSION "cmake-debug")
#	message("cannot get git version!")
#else()
#	message("build version: ${GIT_VERSION}")
#endif()
#
#set(GO_BUILD_FLAGS_DEBUG -X main.Version=${GIT_VERSION})
set(GO_BUILD_FLAGS_RELEASE -s -w)

if (DEFINED CMAKE_BUILD_TYPE)
	string(TOUPPER ${CMAKE_BUILD_TYPE} UPPER_BUILD_TYPE)
endif()

if(WIN32)
	if(NOT DEFINED ENV{SDL2_ROOT})
		message(FATAL_ERROR, "please set the SDL2 root path to envrionment variable SDL2_ROOT (to include, bin, ...)")
	endif()
	if(NOT DEFINED ENV{MINGW_ROOT})
		message(FATAL_ERROR, "please set the mingw64 root path to envrionment variable MINGW_ROOT")
	endif()
endif()

if(LINUX)
	if(NOT DEFINED TARGET_TAG)
		set(TARGET_TAG "linux")
	endif()
	if (${TARGET_TAG} STREQUAL "rpi")
		set(GO_C_INCLUDE_DIRS "$ENV{SYSROOT}/usr/lib/aarch64-poky-linux/14.2.0/include -I$ENV{SYSROOT}/usr/include -I${CMAKE_CURRENT_LIST_DIR}/disp/epaper")
		set(GO_C_LIBRARIES_FORMATTED "-L$ENV{SYSROOT}/usr/lib -L${CMAKE_CURRENT_BINARY_DIR} -llgpio -lrgpio -leplib")
	else()
		find_package(SDL2 REQUIRED NO_MODULE_VERSION)
		set(GO_C_INCLUDE_DIRS ${SDL2_INCLUDE_DIR})
		set(GO_C_LIBRARY_DIRS ${SDL2_LIBRARY_DIRS})
		set(GO_C_LIBRARIES ${SDL2_LIBRARIES})
		#string(REPLACE ";" " -L" GO_C_LIBRARY_DIRS_FORMATTED "-L${GO_C_LIBRARY_DIRS}")
		string(REPLACE ";" " -l" GO_C_LIBRARIES_TEMP " -l${GO_C_LIBRARIES}")
		#string(REPLACE "-l-l" "-l" GO_C_LIBRARIES_FORMATTED "${GO_C_LIBRARIES_TEMP}")
	endif()
elseif(WIN32)
	set(GO_C_INCLUDE_DIRS $ENV{SDL2_ROOT}/include)
	cmake_path(CONVERT ${GO_C_INCLUDE_DIRS} TO_CMAKE_PATH_LIST GO_C_INCLUDE_DIRS)
	set(GO_C_LIBRARY_DIRS $ENV{SDL2_ROOT}/lib)
	cmake_path(CONVERT ${GO_C_LIBRARY_DIRS} TO_CMAKE_PATH_LIST GO_C_LIBRARY_DIRS)
	set(GO_C_LIBRARIES "-lSDL2 -lSDL2main -lmingw32 -mwindows")
	if(NOT DEFINED TARGET_TAG)
		set(TARGET_TAG "windows")
	endif()
	string(REPLACE ";" " -L" GO_C_LIBRARY_DIRS_FORMATTED "-L${GO_C_LIBRARY_DIRS}")
	string(REPLACE ";" " -l" GO_C_LIBRARIES_FORMATTED " ${GO_C_LIBRARIES}")
endif()

string(REPLACE ";" "  -I" GO_C_CFLAGS_FORMATTED " -I${GO_C_INCLUDE_DIRS}")
set(GO_SDL_CFLAGS "${GO_C_CFLAGS_FORMATTED}")
set(GO_SDL_LDFLAGS "${GO_C_LIBRARY_DIRS_FORMATTED}${GO_C_LIBRARIES_FORMATTED}")

if(LINUX)
	set(SET_CGO_ENV_FILE "${CMAKE_CURRENT_BINARY_DIR}/env.sh")
    configure_file(
        "${CMAKE_CURRENT_SOURCE_DIR}/disp/env.sh.in"
        "${SET_CGO_ENV_FILE}" @ONLY
    )

	set(RUN_CGO_ENV_FILE . ${SET_CGO_ENV_FILE})
	add_custom_command(
		OUTPUT ${SET_CGO_ENV_FILE}
		COMMAND ${CMAKE_COMMAND} -E chmod +x "${SET_CGO_ENV_FILE}"
	)
	add_custom_target(
		init_env ALL DEPENDS ${SET_CGO_ENV_FILE}
	)
else()
	set(SET_CGO_ENV_FILE "${CMAKE_CURRENT_BINARY_DIR}/env.bat")
    configure_file(
        "${CMAKE_CURRENT_SOURCE_DIR}/disp/env.bat.in"
        "${SET_CGO_ENV_FILE}" @ONLY
    )

	set(RUN_CGO_ENV_FILE call ${SET_CGO_ENV_FILE})
endif()

function(build_go NAME SRC_PATH TAGS EXTRA_DEPEND GO_DEPEND)
	add_custom_command(
		OUTPUT .${NAME}.gocmd
		DEPENDS ${GO_DEPEND}
		COMMAND ${RUN_CGO_ENV_FILE}
		COMMAND cd ${SRC_PATH} && go mod tidy
		COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_CURRENT_BINARY_DIR}/output
		COMMAND ${CMAKE_COMMAND} -E echo go build -C ${SRC_PATH} -tags ${TAGS} -x -o "${CMAKE_CURRENT_BINARY_DIR}/output/${NAME}" -ldflags="${GO_BUILD_FLAGS_${UPPER_BUILD_TYPE}}"
		COMMAND ${RUN_CGO_ENV_FILE} && go build -C ${SRC_PATH} -tags ${TAGS} -o "${CMAKE_CURRENT_BINARY_DIR}/output/${NAME}" -ldflags="${GO_BUILD_FLAGS_${UPPER_BUILD_TYPE}}"
		WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR}
		COMMENT "Building go project ... ${NAME}"
	)

	add_custom_target(
		${NAME}
		ALL DEPENDS .${NAME}.gocmd ${EXTRA_DEPEND}
	)
endfunction()

if (${TARGET_TAG} STREQUAL "linux")
	#build_go(cal_test "${CMAKE_CURRENT_LIST_DIR}" "cal_test" "${SET_CGO_ENV_FILE}")
	build_go(disp_test "${CMAKE_CURRENT_LIST_DIR}" "disp_test" ${SET_CGO_ENV_FILE} "init_env")
	build_go(ggcal "${CMAKE_CURRENT_LIST_DIR}" "linux" ${SET_CGO_ENV_FILE} "init_env")
elseif (${TARGET_TAG} STREQUAL "windows")
	#build_go(cal_test.exe "${CMAKE_CURRENT_LIST_DIR}" "cal_test" "${SET_CGO_ENV_FILE}")
	build_go(disp_test.exe "${CMAKE_CURRENT_LIST_DIR}" "disp_test" "${SET_CGO_ENV_FILE}" "init_env")
	build_go(ggcal.exe "${CMAKE_CURRENT_LIST_DIR}" "windows" "${SET_CGO_ENV_FILE}" "init_env")
	file(COPY $ENV{SDL2_ROOT}/bin/SDL2.dll DESTINATION ${CMAKE_CURRENT_BINARY_DIR}/output)
elseif (${TARGET_TAG} STREQUAL "rpi")
	build_go(ggcal "${CMAKE_CURRENT_LIST_DIR}" "rpi" ${SET_CGO_ENV_FILE} eplib)
endif()

if (${TARGET_TAG} STREQUAL "rpi")
	add_library(eplib STATIC 
		${CMAKE_CURRENT_LIST_DIR}/disp/epaper/DEV_Config.c
		${CMAKE_CURRENT_LIST_DIR}/disp/epaper/EPD_12in48b.c
	)
	target_include_directories(eplib PUBLIC
		${CMAKE_CURRENT_LIST_DIR}/disp/epaper
		$ENV{SYSROOT}/usr/include
	)
endif()