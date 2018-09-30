FROM golang:1.9.4

RUN go get github.com/jekabolt/ARGraffti-back && \
    rm -rf $GOPATH/src/github.com/jekabolt/ARGraffti-back && \
    cd $GOPATH/src/github.com/jekabolt && \ 
    git clone https://github.com/jekabolt/ARGraffti-back.git && \ 
    cd ARGraffti-back 

RUN go get github.com/jekabolt/config

RUN cd $GOPATH/src/github.com/jekabolt/ARGraffti-back && \
    make build 


WORKDIR $GOPATH/src/github.com/Appscrunch/Multy-BTC-node-service/cmd

RUN echo "VERSION 02"

ENTRYPOINT $GOPATH/src/github.com/Appscrunch/Multy-BTC-node-service/cmd/kandinsky
