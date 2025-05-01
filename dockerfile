FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o server .

FROM alpine

WORKDIR /app
COPY --from=builder /app/server .
COPY client/ client/

RUN mkdir -p data

EXPOSE 3030

CMD ["./server"]