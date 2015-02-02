IMAGE = mlafeldt/chef-runner:rpm

all: image

image:
	docker build --force-rm -t $(IMAGE) $(CURDIR)

.PHONY: all image
