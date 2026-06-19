APP_NAME=subscriptions-service
MAIN=./cmd/app
SWAG=C:\Users\Admin\go\bin\swag.exe

run:
go run $(MAIN)

test:
go test ./...

fmt:
gofmt -w .

tidy:
go mod tidy

swagger:
$(SWAG) init -g cmd/app/main.go

docker-up:
docker-compose up --build

docker-down:
docker-compose down

logs:
docker logs subscriptions_app --tail 100
