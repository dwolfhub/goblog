
create-migration:
	@read -p "Enter Migration Name: " name; \
	migrate create -ext sql -dir migrations $$name

test:
	go test ./...