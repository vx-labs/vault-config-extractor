FROM vxlabs/dep as builder

WORKDIR $GOPATH/src/github.com/vx-labs/vault-config-extractor
COPY Gopkg* ./
RUN dep ensure -vendor-only
COPY . ./
RUN go test ./... && \
    go build -buildmode=exe -a -o /bin/vault-config-extractor ./main.go

FROM alpine
COPY --from=builder /bin/vault-config-extractor /bin/vault-config-extractor

