.PHONY: gen-docs
gen-docs:
	@swag init -g cmd/api/main.go -d . && swag fmt