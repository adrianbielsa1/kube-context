FROM golang:latest

WORKDIR /app/data

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go .

RUN mkdir icons

COPY icons/kubernetes.png icons/kubernetes.png

# installing requirements to get and extract prebuilt binaries
RUN apt-get update && apt-get install -y \
    xz-utils \
    curl \
    && rm -rf /var/lib/apt/lists/*#Getting prebuilt binary from llvm 
RUN curl -SL https://github.com/llvm/llvm-project/releases/download/llvmorg-10.0.0/clang+llvm-10.0.0-x86_64-linux-gnu-ubuntu-18.04.tar.xz \
    | tar -xJC . && \
    mv clang+llvm-10.0.0-x86_64-linux-gnu-ubuntu-18.04 clang_10 && \
    echo ‘export PATH=/clang_10/bin:$PATH’ >> ~/.bashrc && \
    echo ‘export LD_LIBRARY_PATH=/clang_10/lib:$LD_LIBRARY_PATH’ >> ~/.bashrc

RUN CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build .

RUN ls
