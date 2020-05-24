run-sample: build-sample
	@ \
	bash -c 'out/sample-skill 2> >(jq)'

build-sample:
	@ \
	go build -o out/sample-skill sample/cmd/main.go

