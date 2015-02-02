IMAGE = mlafeldt/chef-runner:deb

all: image

image:
	docker build --force-rm -t $(IMAGE) $(CURDIR)

.PHONY: all image
