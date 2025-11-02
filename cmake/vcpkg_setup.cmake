if (APPLE)
    if (NOT DEFINED CMAKE_OSX_ARCHITECTURES)
        execute_process(COMMAND "uname" "-m" OUTPUT_VARIABLE UNAME_RESULT OUTPUT_STRIP_TRAILING_WHITESPACE)
        set(CMAKE_OSX_ARCHITECTURES ${UNAME_RESULT} CACHE STRING "osx_architectures")
    endif()

    if (CMAKE_OSX_ARCHITECTURES STREQUAL "arm64")
        set(CMAKE_OSX_DEPLOYMENT_TARGET 11.0)
        message(STATUS "Building for Apple Silicon Mac computers")
        set(VCPKG_TARGET_TRIPLET arm64-osx-min-11-0)
    elseif (CMAKE_OSX_ARCHITECTURES STREQUAL "x86_64")
        set(CMAKE_OSX_DEPLOYMENT_TARGET 10.15)
        message(STATUS "Building for Intel based Mac computers")
        set(VCPKG_TARGET_TRIPLET x64-osx-min-10-15)
    else ()
        message(FATAL_ERROR "Unknown value for CMAKE_OSX_ARCHITECTURE. Please use one of \"arm64\" and \"x86_64\". Multiple architectures are not supported.")
    endif ()
endif()

set(_vcpkg_toolchain "${CMAKE_SOURCE_DIR}/vcpkg/scripts/buildsystems/vcpkg.cmake")

# Ensure the vcpkg submodule is available so that CodeQL and other automated
# environments that perform a shallow checkout can still configure the project.
if (NOT EXISTS "${_vcpkg_toolchain}")
    execute_process(
        COMMAND git submodule update --init --depth 1 --recursive vcpkg
        WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
        RESULT_VARIABLE _vcpkg_submodule_result
        ERROR_QUIET
    )

    if (NOT _vcpkg_submodule_result EQUAL 0)
        message(WARNING "Failed to initialise vcpkg submodule (exit code ${_vcpkg_submodule_result}).")
    endif ()
endif ()

if (EXISTS "${_vcpkg_toolchain}")
    set(CMAKE_TOOLCHAIN_FILE "${_vcpkg_toolchain}")
else ()
    message(WARNING "vcpkg toolchain file not found; continuing without vcpkg integration.")
endif ()
