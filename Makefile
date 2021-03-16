PROVIDER_ONLY_PKGS=$(shell go list ./... | grep -v "/vendor/" | grep -v "tools")

default: build

build:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	go build -o terraform-provider-mongoatlas .

deploy:
	curl -u${username}:${password} -T terraform-provider-mongoatlas "http://artifactory.yoox.net/artifactory/ecops-fe/terraform-providers/terraform-provider-mongoatlas"

test:
	TF_ACC=1 go test -v 

plan:
	@terraform plan

clean:
	rm terraform-provider-mongoatlas
