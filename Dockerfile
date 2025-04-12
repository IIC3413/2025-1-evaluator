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
# Create submissions user and group.
RUN groupadd --gid 1001 subgroup
RUN useradd subuser --uid 1001 --gid 1001 --shell /bin/bash
ENTRYPOINT ["./main"]
