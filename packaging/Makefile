CHEF_RUNNER_IMAGE ?= mlafeldt/chef-runner
PACKAGECLOUD_REPO ?= mlafeldt/chef-runner

DEB_DISTROS = ubuntu/precise ubuntu/trusty ubuntu/utopic \
	      debian/squeeze debian/wheezy debian/jessie
RPM_DISTROS = el/6 el/7

all: build

deb_image:
	docker build --force-rm -t $(CHEF_RUNNER_IMAGE):deb $(CURDIR)/deb

rpm_image:
	docker build --force-rm -t $(CHEF_RUNNER_IMAGE):rpm $(CURDIR)/rpm

images: deb_image rpm_image

deb_build: deb_image
	docker run -it --rm -v $(CURDIR):/data $(CHEF_RUNNER_IMAGE):deb

rpm_build: rpm_image
	docker run -it --rm -v $(CURDIR):/data $(CHEF_RUNNER_IMAGE):rpm

build: deb_build rpm_build

deb_push: deb_build
	@for distro in $(DEB_DISTROS); do \
		pkgcloud-push $(PACKAGECLOUD_REPO)/$$distro pkg/*.deb || exit 1; \
	done

rpm_push: rpm_build
	@for distro in $(RPM_DISTROS); do \
		pkgcloud-push $(PACKAGECLOUD_REPO)/$$distro pkg/*.rpm || exit 1; \
	done

push: deb_push rpm_push

release: build push

clean:
	$(RM) -r cache tmp-build tmp-dest

clobber: clean
	$(RM) -r pkg

.PHONY: all deb_image rpm_image images deb_build rpm_build build \
	deb_push rpm_push push release clean clobber
