all: test

bootstrap:
	@script/bootstrap

generate:
	@go generate -x ./...

update_omnibus:
	@curl -L https://www.chef.io/chef/install.sh >chef/omnibus/assets/install.sh
	@go generate -x ./chef/omnibus

lint:
	@script/lint

test:
	@script/test

build:
	@script/build

release:
	@script/build --release

packages:
	$(MAKE) -C packaging build

clean:
	$(RM) -r build

.PHONY: all bootstrap generate update_omnibus \
	lint test build release packages clean
