.PHONY: build run install clean docker run_svc run_docker_svc test tag_v1 deploy

HOSTING_SERVER=fhdemo

all: rebuild # builds so fast that full rebuils are practical.

rebuild: clean build

build: bin/fh-geo-svc
	@echo "Build OK"
	
bin/fh-geo-svc:
	$(eval GOPATH=$(shell pwd))
	CGO_ENABLED=0 GOPATH=$(GOPATH) go build -o $@ fh-geo-svc

unittest:
	$(eval GOPATH=$(shell pwd))
	CGO_ENABLED=0 GOPATH=$(GOPATH) go test fh-libgeo
	
docker:
	docker build -t fh-geo-svc:v1-dev .

clean:
	rm -rf ./bin/* docker/fh-geo-svc/artifacts/*

run: run_svc

parse: rebuild
	cd db && ./recreate.sh --force
	LOG_LEVEL=3 ./bin/./fh-geo-svc parse
	
run_svc:
	./bin/./fh-geo-svc

run_docker_svc: docker/fh-geo-svc
	docker run --rm -ti fh-geo-svc:v1-dev
	
test:
	make rebuild
	#cd db && ./recreate.sh --force
	#DEBUG_REQUEST_BODY=1 LOG_LEVEL=3 ./bin/./fh-geo-svc
	
tag_v1:
	docker tag fh-geo-svc:v1-dev fh-geo-svc:v1
	
deploy: docker test tag_v1
	docker save --output ./bin/fh-geo-svc.tar.gz fh-geo-svc:v1
	scp ./bin/fh-geo-svc.tar.gz $(HOSTING_SERVER):/opt/fh-geo-svc.tar.gz
	ssh $(HOSTING_SERVER) systemctl restart fh-geo-svc
	#ssh $(HOSTING_SERVER) 
	
deploy-data:
	rsync -avzr data/ $(HOSTING_SERVER):/opt/data/
	
