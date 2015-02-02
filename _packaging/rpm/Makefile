#
# Build Debian package of chef-runner and push it to packagecloud.io for all
# supported Debian/Ubuntu distributions.
#

CHEF_RUNNER_IMAGE = mlafeldt/chef-runner:centos
PACKAGECLOUD_REPO = mlafeldt/chef-runner

DISTROS = el/6 el/7

all: build

image:
	docker build --force-rm -t $(CHEF_RUNNER_IMAGE) $(CURDIR)

build: image
	docker run -it --rm -v $(CURDIR):/data $(CHEF_RUNNER_IMAGE)

push:
	@for distro in $(DISTROS); do \
		pkgcloud-push $(PACKAGECLOUD_REPO)/$$distro pkg/*.rpm || exit 1; \
	done

release: build push

clean:
	$(RM) -r cache tmp-build tmp-dest

clobber: clean
	$(RM) -r pkg

.PHONY: all image build push release clean clobber
