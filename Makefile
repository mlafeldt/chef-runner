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

deb:
	$(MAKE) -C _packaging/deb build

rpm:
	$(MAKE) -C _packaging/rpm build

clean:
	$(RM) -r .cover build

.PHONY: all bootstrap generate update_omnibus \
	lint test coverage build release deb rpm clean
