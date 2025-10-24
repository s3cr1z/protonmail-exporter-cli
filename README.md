# Proton Mail Export

Proton Mail Export allows you to export your emails as eml files.

## Features

- Export all emails or filter by specific folders/labels
- Support for backup and restore operations
- Cross-platform support (Linux, Mac, Windows)

For information on using the folder filter feature, see [FILTER_EXPORT_USAGE.md](FILTER_EXPORT_USAGE.md).

# Building

## Requirements

- C++ 17 compatible compiler
  - GCC/Clang (Linux/Mac)
  - MSVC 2022 (Windows)
- CMake >= 3.23
- Go >= 1.24

## Fetch submodules

```
git submodule update --init --recursive
```

## Linux/Mac

```
cmake -S. -B $BUILD_DIR -G <Insert favorite Generator>
cmake --build $BUILD_DIR
```

## Windows

```
cmake -S. -B $BUILD_DIR -G "Visual Studio 17 2022" -DVCPKG_TARGET_TRIPLET=x64-windows-static
cmake --build $BUILD_DIR --config Release
```

**Note:** An active internet connection is required in otder to dowload a standalone MingW compiler in order to compile
the CGO module.

## Layout

- [go-lib](go-lib): CGO Shared library implementation
- [lib](lib): C++ shared library over the exported C interface from [go-lib](go-lib)
- [cli](cli): CLI application
