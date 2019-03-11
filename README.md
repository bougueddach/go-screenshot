# go-screenshot

This a CLI application that allow you to take web pages screenshots using the command interface


# Clone it

    `mkdir -p $GOPATH/src/ && cd $_ && git clone https://github.com/bougueddach/go-screenshot.git`

# Build it!

Here are the steps to get you started on the project:
1) Install dependencies
    `dep ensure`
2) Install go-screenshot command
    `make install`

3) Usage Examples
   <br /> `go-screenshot url --values "https://detectify.com"`
   <br /> `go-screenshot url --values "https://medium.com;https://detectify.com"`
