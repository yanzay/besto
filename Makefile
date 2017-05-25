build:
	go build -v -i .
dev: build
	./besto --local --log-level trace
deploy:
	ssh root@yanzay.com "cd infra; docker-compose pull besto; docker-compose up -d"
clean:
	rm besto.db
	rm besto
