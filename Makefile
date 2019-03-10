APPNAME := go-screenshot

clean:
	@rm -f ./$(APPNAME)

build: clean
	go build -v -o $(APPNAME) cmd/main.go

run: build
	./$(APPNAME)




