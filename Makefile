.PHONY: build-and-push

tag = citystatetz

build:
	docker build -t $(tag) .

