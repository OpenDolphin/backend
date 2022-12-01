FROM golang:1.19-alpine AS builder
COPY . /app
WORKDIR /app
RUN apk add --no-cache make
RUN ls -la ./cmd
RUN make build

FROM alpine:3.16
COPY --from=builder /app/build/opendolphin-backend /usr/bin/opendolphin-backend
ENV LISTEN_ADDR=0.0.0.0:5000
EXPOSE 5000
ENTRYPOINT [ "/usr/bin/opendolphin-backend" ]
