FROM golang:alpine AS build

WORKDIR /go/src/app

COPY . .

RUN go build -o main .

FROM alpine

WORKDIR /app

COPY --from=build /go/src/app .

COPY .env /app

EXPOSE 10002

ENTRYPOINT [ "./main" ]