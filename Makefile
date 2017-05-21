build:
	go build -v -i .
dev: build
	./besto --local --log-level trace
linuxbuild:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i -v .
docker: linuxbuild
	docker build -t yanzay/besto .
push: docker
	docker push yanzay/besto
deploy: push
	ssh root@yanzay.com "cd infra; docker-compose pull besto; docker-compose up -d"
clean:
	rm besto.db
	rm besto
