![Intel oneAPI X GoLang](https://raw.githubusercontent.com/raj96/image_filter_oneapi_go/main/ui/assets/logo.png)

# image_filter_oneapi_go
Couple of image filters implemented in Go using routines provided by oneAPI

## Tech stack:
#### UI:
 - Vanilla JS/HTML/CSS
 - Materialize CSS
#### Backend:
 - Go
 - net/http
 - Intel oneAPI MKL

## Build process:
The location to `go` binary is set to `/opt/go/bin/go` using the variable `GO` in `Makefile`, one should change it, if `go` toolchain is installed elsewhere.

The build process assumes that the Intel oneApi toolkit is installed under `/opt/intel/oneapi`, one should change the variable `ONEAPI_ROOT` to the appropriate value in `Makefile`, in case the oneApi toolkit is installed somewhere else.
Run `make all` and check if `cgo` is able to detect all the libraries/header files properly. The `#cgo` flags are used inside `api/grayscale.go` and `api/gaussian_blur.go` (basically all the APIs that perform convoultion will use Intel's oneAPI MKL routines for performance boost).

If `make all` is successful. One can go ahead and run `make docker-image`, BUT make sure to change the docker image name inside `Makefile` if the plan is to make changes to the source code and push it on Docker Hub.

If `libtbb.so` is not found, make sure the `ONEAPI_ROOT` variable is set properly.

## Port configuration
The default port is set to 8080, but can be overridden by setting the environment variable `PORT`.

## oneAPI MKL routines used
 - [vsMul](https://www.intel.com/content/www/us/en/docs/onemkl/developer-reference-c/2023-2/v-mul.html)
 - [cblas_sasum](https://www.intel.com/content/www/us/en/docs/onemkl/developer-reference-c/2023-0/cblas-asum.html)
