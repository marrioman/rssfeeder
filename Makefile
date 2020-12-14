local-pg:
	docker run -d --name test-feeder-postgres -p 5444:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=feederdb postgres:11-alpine

clean-local-pg:
	docker rm -f test-feeder-postgres