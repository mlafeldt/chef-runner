generate:
	@go generate -x ./...

update_omnibus:
	@curl -Ls https://www.opscode.com/chef/install.sh >chef/omnibus/assets/install.sh
	@go generate -x ./chef/omnibus

.PHONY: generate update_omnibus
