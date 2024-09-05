cmake_minimum_required(VERSION 3.1)
project(rnnoise)

option(RNNOISE_COMPILE_OPUS ON)

if(RNNOISE_COMPILE_OPUS)
	add_definitions(-DCOMPILE_OPUS)
endif()

# Ignore CRT warnings on MSVC
if(MSVC)
  add_definitions(-D_CRT_SECURE_NO_WARNINGS)
endif()

# Get source files
file(GLOB SOURCES "src/*.c"  "src/*.h" "include/*.h")

# Build rnnoise
add_definitions(-DRNNOISE_BUILD)

# Compile the library
add_library(rnnoise ${SOURCES})

# Build DLL if needed
if(BUILD_SHARED_LIBS)
  if(WIN32)
    target_compile_definitions(rnnoise PRIVATE DLL_EXPORT)
  else()
    include(CheckCCompilerFlag)
    check_c_compiler_flag(-fvisibility=hidden COMPILER_HAS_HIDDEN_VISIBILITY)
    if(COMPILER_HAS_HIDDEN_VISIBILITY)
      set_target_properties(rnnoise PROPERTIES C_VISIBILITY_PRESET hidden)
    endif()
  endif()
endif()

include(GNUInstallDirs)

# Include dirs
target_include_directories(rnnoise PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
        PRIVATE src)

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