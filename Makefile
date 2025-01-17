EXECUTABLE=jira-to-pdf
PACKAGE=github.com/prsolucoes
LOG_FILE=/var/log/${EXECUTABLE}.log
GOFMT=gofmt -w
GODEPS=go get -u

.DEFAULT_GOAL := help

# general
help:
	@echo "Type: make [rule]. Available options are:"
	@echo ""
	@echo "- help"
	@echo "- build"
	@echo "- install"
	@echo "- format"
	@echo "- deps"
	@echo "- start"
	@echo "- stop"
	@echo "- update"
	@echo "- build-all"
	@echo ""

build:
	go build -o ${EXECUTABLE}

install:
	go install

format:
	${GOFMT} main.go

deps:
	${GODEPS} github.com/andygrunwald/go-jira
	${GODEPS} github.com/jung-kurt/gofpdf
	${GODEPS} golang.org/x/text/encoding/charmap

stop:
	pkill -f ${EXECUTABLE}

start:
	-make stop
	cd ${GOPATH}/src/${PACKAGE}/${EXECUTABLE}
	nohup ${EXECUTABLE} >> ${LOG_FILE} 2>&1 </dev/null &

update:
	git pull origin master
	make install

build-all:
	rm -rf build

	mkdir -p build/linux32
	env GOOS=linux GOARCH=386 go build -o build/linux32/${EXECUTABLE} -v ${PACKAGE}/${EXECUTABLE}

	mkdir -p build/linux64
	env GOOS=linux GOARCH=amd64 go build -o build/linux64/${EXECUTABLE} -v ${PACKAGE}/${EXECUTABLE}

	mkdir -p build/darwin64
	env GOOS=darwin GOARCH=amd64 go build -o build/darwin64/${EXECUTABLE} -v ${PACKAGE}/${EXECUTABLE}

	mkdir -p build/windows32
	env GOOS=windows GOARCH=386 go build -o build/windows32/${EXECUTABLE}.exe -v ${PACKAGE}/${EXECUTABLE}

	mkdir -p build/windows64
	env GOOS=windows GOARCH=amd64 go build -o build/windows64/${EXECUTABLE}.exe -v ${PACKAGE}/${EXECUTABLE}