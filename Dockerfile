# Command to buiild: sudo docker build --no-cache -t uniiotx-querier .
FROM golang:1.22.5-alpine3.20 as build
WORKDIR /go/src/github.com/Bedrock-Technology/uniiotx-querier
COPY . .
RUN ls -la
RUN go get -d ./...
RUN apk --update add build-base && cd app/ && GOOS=linux go build -o uniiotx-querier .

FROM alpine:3.20
WORKDIR /app/
COPY --from=build /go/src/github.com/Bedrock-Technology/uniiotx-querier/app/uniiotx-querier ./
CMD ["./uniiotx-querier"]
