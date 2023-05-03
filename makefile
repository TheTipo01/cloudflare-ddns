build:
	go build -o out/main main.go
	go build -buildmode=plugin -o out/ipify.so plugins/ipify/ipify.go
	go build -buildmode=plugin -o out/skywifi.so plugins/skywifi/skywifi.go
	go build -buildmode=plugin -o out/vodafone.so plugins/vodafone/vodafone.go
	go build -buildmode=plugin -o out/openwrt.so plugins/openwrt/openwrt.go
	cp -R config out