
init:
	@echo "Create reports folder"
	mkdir -p reports

test: init
	@echo "Unit test is starting..."
	go test -v ./...
	go test -coverprofile ./reports/cover-report.out ./...
	go tool cover -html ./reports/cover-report.out
	@echo "Unit test done"
