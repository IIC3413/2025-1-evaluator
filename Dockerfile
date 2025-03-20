FROM golang:1.24.1
RUN apt-get update && apt-get -y install g++ cmake
WORKDIR /app
# Copy dependency files.
COPY go.sum go.mod ./
RUN go mod download
# Copy compile files.
COPY internal/ internal/
COPY main.go ./
RUN go build -o main .
# Copy IO
COPY wkdir/ wkdir/
ENTRYPOINT ["./main", "-c=config/config.yaml"]
