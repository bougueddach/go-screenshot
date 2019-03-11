# go-screenshot

This a CLI application that allow you to take web pages screenshots using the command line interface, it is inspired from gowitness :wink:.


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

# Scalability thoughts!
I don't have any experience scaling apps but here are my thoughts about what we need:

    * Add a metrics implementation for monitoring
    * Probably we will need a load balancer
    * Probably we will need a queuing systeme depending on the use case that we want
    * If we want to continiously access the requests historic I suggest to use a real db to benefit from indexing
    