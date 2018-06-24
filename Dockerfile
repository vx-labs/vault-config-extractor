FROM vxlabs/glide as builder

WORKDIR $GOPATH/src/github.com/vx-labs/vault-config-extractor
COPY glide* ./
RUN glide install -v
COPY . ./
RUN go test $(glide nv) && \
    go build -buildmode=exe -a -o /bin/vault-config-extractor ./main.go

FROM alpine
COPY --from=builder /bin/vault-config-extractor /bin/vault-config-extractor

