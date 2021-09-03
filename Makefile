.PHONY: all clean build

APP_NAME=gui-mini-ffmpeg

all: prod

prod: clean compile

## default run
run:
	cd src; go run . $(ARGS)

## check race condition
race:
	cd src; go run -race .

## default build
build:	
	cd src; cp ./main.glade ../build/main.glade; go build -o ../build/${APP_NAME} .

## production build (strip the debugging information)
compile:
	cd src; cp ./main.glade ../build/main.glade; GOOS=linux GOARCH=amd64 go build -o ../build/${APP_NAME} .
##	cd src;	GOOS=windows GOARCH=amd64 go build -o ../build/${APP_NAME}.exe .

## clear cache and remove builds
clean:
	go clean
	rm -f build/${APP_NAME}
	rm -f build/main.glade
