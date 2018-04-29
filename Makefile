.PHONY: test
test:
	docker build --rm -t goenchant-test .
	docker run --rm -ti goenchant-test
