# Proton Mail Export Repository Summary

## 1. Repository Purpose
This repository contains the source code for Proton Mail Export, a command-line tool that allows users to export their emails from Proton Mail as EML files. It provides powerful filtering options to export specific subsets of emails based on various criteria.

## 2. Repository Setup
The project is a C++ and Go hybrid, using CMake for the build system. It has a client-server architecture where the C++ `cli` communicates with the `go-lib` for email filtering and exporting. For a more detailed guide on building the project, refer to the [official documentation](https-placeholder).

## 3. Repository Structure
The repository is organized into the following key directories:
- `go-lib`: Contains the core Go library responsible for handling email filtering and exporting.
- `lib`: A C++ shared library that wraps the Go library, providing a C-compatible interface.
- `cli`: A command-line interface (CLI) application that allows users to interact with the export tool.
- `ci`: Contains CI/CD pipeline configurations for linting, building, and deploying the project.
- `cmake`: Includes CMake scripts and modules for managing the build process.

## 4. CI Checks
The CI pipeline is defined in `.gitlab-ci.yml` and includes the following checks:
- **Linting**:
  - `go-lib-lint`: Runs `golangci-lint` to check for issues in the Go code.
  - `clang-format-check`: Ensures the C++ code adheres to the defined coding style.
- **Building**:
  - Compiles the project for various platforms, including Linux, macOS, and Windows.
- **Deployment**:
  - Deploys the application when a new tag is pushed to the repository.
