all: test

bootstrap:
	@script/bootstrap

generate:
	@go generate -x ./...

update_omnibus:
	@curl -Ls https://www.opscode.com/chef/install.sh >chef/omnibus/assets/install.sh
	@go generate -x ./chef/omnibus

lint:
	@script/lint

test:
	@script/test

coverage:
	@script/coverage --html

build:
	@script/build

release:
	@script/build --release

packages:
	$(MAKE) -C packaging build

clean:
	$(RM) -r .cover build

.PHONY: all bootstrap generate update_omnibus \
	lint test coverage build release packages clean
