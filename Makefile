checkout:
	git clone https://github.com/ipfs/testground.git
	cd testground && git checkout jim/cloud-run-experiment
	mkdir -p targets
	cd targets && git clone https://github.com/ipfs/go-ipfs.git

build-targets:
	cd targets/go-ipfs && make build

distclean:
	rm -rf testground targets

run:
	go run main.go

docker-build:
	docker build -t jimpick/testground-web .

docker-run:
	docker run -p 8099:8099 --name testground jimpick/testground-web

docker-rm:
	docker rm testground

gcloud-build:
	go mod tidy
	go mod vendor
	gcloud builds submit --tag gcr.io/testground/testground-web

gcloud-deploy:
	gcloud beta run deploy testground-web --image gcr.io/testground/testground-web --platform managed --region europe-west1 --timeout 30s

gcloud-list-services:
	gcloud beta run services list --platform managed

deploy: gcloud-build
	make gcloud-deploy
