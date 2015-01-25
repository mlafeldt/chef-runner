#
# Build Debian package of chef-runner and push it to packagecloud.io for all
# supported Debian/Ubuntu distributions. All using Docker containers.
#

CHEF_RUNNER_IMAGE  = mlafeldt/chef-runner:ubuntu
PACKAGECLOUD_IMAGE = mlafeldt/packagecloud
PACKAGECLOUD_REPO  = mlafeldt/chef-runner

DISTROS = ubuntu/precise ubuntu/trusty ubuntu/utopic \
          debian/squeeze debian/wheezy debian/jessie

all: build

image:
	docker build --force-rm -t $(CHEF_RUNNER_IMAGE) $(CURDIR)

build: image
	docker run -it --rm -v $(CURDIR):/data $(CHEF_RUNNER_IMAGE)

push:
	@for distro in $(DISTROS); do \
		echo "Pushing package for $$distro..."; \
		docker run -it --rm -v $(CURDIR)/pkg:/pkg:ro \
			-e PACKAGECLOUD_TOKEN=$(PACKAGECLOUD_TOKEN) $(PACKAGECLOUD_IMAGE) \
			/bin/bash -c "package_cloud push $(PACKAGECLOUD_REPO)/$$distro /pkg/*.deb" || exit 1; \
	done

release: build push

clean:
	$(RM) -r cache/ tmp-build/ tmp-dest/

clobber: clean
	$(RM) -r pkg/

.PHONY: all image build push release clean clobber
