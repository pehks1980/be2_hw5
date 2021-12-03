#!/bin/zsh
echo Clean db
migrate -database "postgres://test:test@localhost:8181/test?sslmode=disable" -path migrations down
migrate -database "postgres://test:test@localhost:8191/test?sslmode=disable" -path migrations down
migrate -database "postgres://test:test@localhost:8181/test?sslmode=disable" -path migrations up
migrate -database "postgres://test:test@localhost:8191/test?sslmode=disable" -path migrations up
echo Done.