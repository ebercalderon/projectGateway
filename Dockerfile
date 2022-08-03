FROM golang:1.18.3-alpine

WORKDIR /erpgateway

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . ./

RUN go build -o /erp_gateway

EXPOSE 8081
CMD [ "/erp_gateway" ]