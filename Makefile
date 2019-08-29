run:
	go run main.go

docker-build:
	docker build -t jimpick/testcloud-web .

docker-run:
	docker run -p 8099:8099 jimpick/testcloud-web 

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
