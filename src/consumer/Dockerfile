FROM golang:1.22.1-alpine3.19 AS build
RUN apk update && apk add --no-cache tzdata
WORKDIR /app
COPY . .

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

FROM alpine
RUN adduser -D -u 1000 appuser
USER appuser
WORKDIR /app
COPY --from=build /app /app
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=America/Sao_Paulo

CMD ["/app/main"]