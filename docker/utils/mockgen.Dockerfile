FROM --platform=linux/arm64 golang:1.20-alpine

RUN go install github.com/golang/mock/mockgen@v1.6.0

CMD ["go", "generate", "./..."]
