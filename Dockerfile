FROM golang:1.27.1 AS builder

WORKDIR /scry

copy . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/bin/scry

FROM scratch

COPY --from=builder /usr/bin/scry /scry

USER nonroot

EXPOSE 22

CMD ["/scry"]
