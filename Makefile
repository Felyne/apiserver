GO_EXECUTABLE ?= go

BINARY = apiserver
DIST = dist

build:
	${GO_EXECUTABLE} build -o ${BINARY}

clean:
	rm -f ${BINARY}
	rm -rf ${DIST}


.PHONY: build clean
