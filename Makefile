yages_version := 0.1.0

.PHONY: build clean container push delete destroy

build :
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.release=$(yages_version)" -o ./srv-yages ./main.go

clean :
	@rm srv-yages

container :
	@docker build --build-arg yversion=$(yages_version) -t quay.io/mhausenblas/yages:$(yages_version) .

push :
	@docker push quay.io/mhausenblas/yages:$(yages_version)

deploy :
	@kubectl create ns grpc-demo
	@kubectl apply -f app.yaml

destroy :
	@kubectl delete ns grpc-demo
