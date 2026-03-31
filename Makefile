.PHONY: test test-node test-python test-php validate ai-report-test

test: test-node test-python test-php

test-node:
	cd dagger && dagger call test --source=../examples/node --language=node

test-python:
	cd dagger && dagger call test --source=../examples/python --language=python

test-php:
	cd dagger && dagger call test --source=../examples/php-symfony --language=php

validate:
	cd dagger && dagger call validate-yaml --yaml-file=../templates/github/ai-report.yml

ai-report-test:
	cd dagger && dagger call ai-report-test --source=../examples/node
