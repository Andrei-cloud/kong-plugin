FROM kong/go-plugin-tool:2.0.4-alpine-latest AS builder

COPY . /app
WORKDIR /app
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o /app/bin/num2word .

FROM kong:2.7.0

COPY --from=builder /app/bin/num2word /usr/local/bin/num2word
