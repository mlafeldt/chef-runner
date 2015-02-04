DOCKER_IMAGE ?= mlafeldt/chef-runner
PACKAGECLOUD_REPO ?= mlafeldt/chef-runner

DEB_DISTROS = ubuntu/precise ubuntu/trusty ubuntu/utopic \
	      debian/squeeze debian/wheezy debian/jessie
RPM_DISTROS = el/6 el/7

all: build

image:
	docker pull -a $(DOCKER_IMAGE)

deb:
	docker run -it --rm -v $(CURDIR):/data $(DOCKER_IMAGE):deb

rpm:
	docker run -it --rm -v $(CURDIR):/data $(DOCKER_IMAGE):rpm

build: deb rpm

push_deb: deb
	@for distro in $(DEB_DISTROS); do \
		pkgcloud-push $(PACKAGECLOUD_REPO)/$$distro pkg/*.deb || exit 1; \
	done

push_rpm: rpm
	@for distro in $(RPM_DISTROS); do \
		pkgcloud-push $(PACKAGECLOUD_REPO)/$$distro pkg/*.rpm || exit 1; \
	done

push: push_deb push_rpm

release: push

clean:
	$(RM) -r cache

clobber: clean
	$(RM) -r pkg

.PHONY: all image deb rpm build push_deb push_rpm push release clean clobber
