
init:
	@echo "Create reports folder"
	mkdir -p reports
	@echo "Done!"

test: init
	@echo "Unit test is starting..."
	go test -v ./...
	go test -coverprofile ./reports/cover-report.out ./...
	go tool cover -html ./reports/cover-report.out
	@echo "Done!"

godoc:
	@echo "Opening godoc..."
	~/go/bin/godoc -http=:6060 -v
	@echo "Done!"

clean:
	@echo "Remove reports..."
	rm -rf ./reports/*
	@echo "Tidy dependence..."
	go mod tidy
	@echo "Done!"