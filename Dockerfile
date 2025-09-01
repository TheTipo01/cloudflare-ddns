FROM --platform=$BUILDPLATFORM golang:alpine AS build

RUN apk add --no-cache build-base

COPY . /cloudflare-ddns

WORKDIR /cloudflare-ddns

RUN CGO_ENABLED=1 go mod download
RUN CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o out/cloudflare-ddns
RUN CGO_ENABLED=1 go build -trimpath -buildmode=plugin -o out/ipify.so plugins/ipify/ipify.go
RUN CGO_ENABLED=1 go build -trimpath -buildmode=plugin -o out/skywifi.so plugins/skywifi/skywifi.go
RUN CGO_ENABLED=1 go build -trimpath -buildmode=plugin -o out/vodafone.so plugins/vodafone/vodafone.go
RUN CGO_ENABLED=1 go build -trimpath -buildmode=plugin -o out/openwrt.so plugins/openwrt/openwrt.go

FROM alpine

COPY --from=build /cloudflare-ddns/out /opt/cloudflare-ddns
WORKDIR /opt/cloudflare-ddns

CMD ["cloudflare-ddns"]
