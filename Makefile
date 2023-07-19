VERSION=0.0.1
BINARY_NAME=signalk-server-go
FREEBOARD_PACKAGE=@signalk/freeboard-sk
FREEBOARD_VERSION=2.1.0
INSTRUMENTPANEL_PACKAGE=@signalk/instrumentpanel
INSTRUMENTPANEL_VERSION=0.24.0
CANBOARD_VERSION=5.0.1

webapps:
	mkdir -p static/${FREEBOARD_PACKAGE} && wget -cq https://registry.npmjs.org/${FREEBOARD_PACKAGE}/-/freeboard-sk-${FREEBOARD_VERSION}.tgz -O -|tar -xz -C static/${FREEBOARD_PACKAGE} package/public --strip-components 2
	mkdir -p static/${INSTRUMENTPANEL_PACKAGE} && wget -cq https://registry.npmjs.org/${INSTRUMENTPANEL_PACKAGE}/-/instrumentpanel-${INSTRUMENTPANEL_VERSION}.tgz -O -|tar -xz -C static/${INSTRUMENTPANEL_PACKAGE} package/public --strip-components 2	
	
build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} -ldflags="-X 'github.com/wdantuma/signalk-server-go/signalkserver.Version=${VERSION}'" main.go

buildarm:
	GOARCH=arm GOOS=linux go build -o ${BINARY_NAME}-arm -ldflags="-X 'github.com/wdantuma/signalk-server-go/signalkserver.Version=${VERSION}'" main.go	

run: build
	./${BINARY_NAME} --mmsi 244810236 --file-source  samples/nemo-n2k.txt

debug: build
	./${BINARY_NAME} --mmsi 244810236 --debug  --file-source  samples/nemo-n2k.txt

clean:
	go clean
	rm ${BINARY_NAME}
