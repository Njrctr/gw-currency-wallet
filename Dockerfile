FROM golang:1.22.2-alpine


RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o gw-currency-wallet ./cmd/main.go

CMD [ "./gw-currency-wallet" ]