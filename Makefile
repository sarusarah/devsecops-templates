.PHONY: test lint contract

test:
	dagger call test --src .

lint:
	dagger call lint --src .

contract:
	dagger call gitlab-contract --src .
