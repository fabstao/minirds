GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=goapp
MAIN=main.go
    
all: build

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v ${MAIN}

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
    
    
# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ${MAIN} 

docker-build:
	docker build -t gcr.io/fabs-cl-02/mrds1 .
