USER_BINARY=userApp
BROKER_BINARY=brokerApp


## up_build: stops docker-compose if running, build all projects and start docker-compose
up_build: build_user build_broker
	@echo "stops all container if running"
	docker-compose down
	@echo "building (when required) and start docker images"
	docker-compose up --build
	@echo "docker images built and started"

## down: stop docker-compose
down:
	@echo "stoping docker-compose"
	docker-compose down
	@echo "done"

## build_user: build the user binary as linux executable file
build_user:
	@echo "building user binary"
	cd ../user-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_BINARY} ./cmd
	@echo "done"

build_broker:
	@echo "building user binary"
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd
	@echo "done"


gen-userproto:
	cd ../user-service && protoc --go_out=user-proto --proto_path=proto proto/*.proto --go-grpc_out=user-proto
	cd ../broker-service && protoc --go_out=user-proto --proto_path=proto proto/*.proto --go-grpc_out=user-proto