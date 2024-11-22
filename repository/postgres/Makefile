mocks:
	mockgen -source=interfaces.go -package=postgres -destination=interface_mock.go

coverage_out:
	go test -coverprofile=coverage.out && go tool cover -html=coverage.out

coverage:
	go test ./... -cover
