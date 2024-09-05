cmake_minimum_required(VERSION 3.1)
project(rnnoise)


# Set the target platform to Android
set(CMAKE_SYSTEM_NAME Android)

# which compilers to use for C and C++
set(CMAKE_C_COMPILER   $ENV{CC})
set(CMAKE_CXX_COMPILER $ENV{CXX})


# Specify the Android platform and architecture
set(ANDROID_ABI armeabi-v7a)  # Specify your desired ABI, e.g., armeabi-v7a, arm64-v8a, x86, etc.

# Set the toolchain file for Android
# set(CMAKE_TOOLCHAIN_FILE  $ENV{NDK_ROOT}/$ENV{NDK_VERSION}/android-cmake/android.toolchain.cmake)

# Specify the Android API level
set(ANDROID_NATIVE_API_LEVEL 21)  # Specify the desired API level

# Include the Android toolchain file
# include(${CMAKE_TOOLCHAIN_FILE})

# Get source files
file(GLOB SOURCES "src/*.c" "src/*.h" "include/*.h")

# Build rnnoise
add_definitions(-DRNNOISE_BUILD)

# Compile the library
add_library(rnnoise ${SOURCES})

# Specify the target ABI for the library
set_target_properties(rnnoise PROPERTIES ANDROID_ABI ${ANDROID_ABI})

# Include dirs
target_include_directories(rnnoise PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
        PRIVATE src)

# Set install directories
include(GNUInstallDirs)

# Install the library
install(TARGETS rnnoise EXPORT rnnoise_target
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
    ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
)

install(DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/include/
    DESTINATION "include")

install(EXPORT rnnoise_target
    FILE rnnoiseConfig.cmake
    DESTINATION lib/cmake/rnnoise
)