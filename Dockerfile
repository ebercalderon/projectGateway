FROM golang:1.18.3-alpine

WORKDIR /projectgateway

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . ./

RUN go build -o /project_gateway

EXPOSE 8081
CMD [ "/project_gateway" ]