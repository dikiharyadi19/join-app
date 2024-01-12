bin-app:
	env GOOS=${OS} GOARCH=${ARCH} go build -o ${NAME}

docker-build:
	docker-compose build

docs:
	@echo "> Generate Swagger Docs"
	@if ! command -v swag &> /dev/null; then go install github.com/swaggo/swag/cmd/swag@latest ; fi
	@swag init -g transport/rest/rest.go --parseDependency true --parseInternal --always-make

start-up:
	docker-compose up --build -d

mock: mock-client mock-usecase

mock-client:
	mockgen -source=./domain/client/open_library.go -destination=./shared/mock/client/open_library_client_mock.go -package client

mock-usecase:
	mockgen -source=./service/library/usecase/library.go -destination=./shared/mock/usecase/library_usecase_mock.go -package usecase

run-test:
	go test service/library/usecase/library.go service/library/usecase/library_test.go -v -cover

migrate-create:
	migrate create -ext ${ext} -dir ${dir} ${filename}

migrate-up:
	migrate -path ${path} -database ${host} -verbose up

migrate-force:
	migrate -path ${path} -database ${host} force ${version}
