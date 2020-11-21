FROM golang:latest

WORKDIR /root/app

COPY . .

# get dependencies
RUN go get -d -v ./...

CMD cp -r /root/app /workspace
