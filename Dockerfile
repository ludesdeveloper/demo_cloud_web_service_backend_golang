# Step 1, membuat binary file
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/mypackage/myapp/
# Tambahan ENV penting, karena bila tidak ditambahkan, binary tidak jalan pada saat running di scratch
ENV CGO_ENABLED=0
COPY . .
# Melakukan instalasi terhadap library yang kita gunakan
RUN go get -d -v
RUN go build -o /go/bin/main
# Step 2, meletakkan binary file diatas scratch
FROM scratch
COPY --from=builder /go/bin/main /go/bin/main
ENTRYPOINT ["/go/bin/main"]