APPNAME := go-screenshot

clean:
	@rm -f ./$(APPNAME)

install:
	go install

build: clean
	go build -v -o $(APPNAME) main.go

run: build
	./$(APPNAME)

