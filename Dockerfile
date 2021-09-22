FROM golang:1.16-buster AS build

RUN go version

COPY . /github.com/Hudayberdyyev/image-service/
WORKDIR /github.com/Hudayberdyyev/image-service/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o file_storage cmd/main.go

FROM alpine:latest

WORKDIR /

COPY --from=build /github.com/Hudayberdyyev/image-service/file_storage .
COPY --from=build /github.com/Hudayberdyyev/image-service/configs configs/

EXPOSE 6670

CMD ["./file_storage"]