test:
	go test -v -failfast
	docker-compose -f e2etest/docker-compose.yaml up -d
	go test -v -failfast ./e2etest
	docker-compose -f e2etest/docker-compose.yaml down