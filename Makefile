ONEAPI_ROOT="/opt/intel/oneapi"
CC="${ONEAPI_ROOT}/compiler/latest/linux/bin/icx"
CXX="${ONEAPI_ROOT}/compiler/latest/linux/bin/icpx"
GO="/opt/go/bin/go"

docker-image: all
	cp "${ONEAPI_ROOT}/tbb/latest/lib/intel64/gcc4.8/libtbb.so.12" .
	docker build -t goxoneapi_image_filter .

all:
	${GO} build .