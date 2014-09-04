## v0.7.0 (unreleased)

* Use run list from JSON file if it contains the `run_list` attribute. Recipes
  passed on the command line will still override this list, and
  `recipes/default.rb` is still the default. (Thanks to @arosenhagen who
  requested this feature.)
* Strip all non-cookbook files after resolving dependencies. This ensures that
  only essential cookbook files are copied to target machines, further saving
  time.
* Extend `script/build` to auto-generate chef-runner's [Homebrew formula] when
  building a new release with `--release`.

[Homebrew formula]: https://github.com/mlafeldt/homebrew-formulas/blob/master/Formula/chef-runner.rb

## v0.6.0 (Aug 28 2014)

FEATURES:

* Add ability to install Chef before provisioning. This allows you to, for
  example, provision bare servers that have nothing installed but the base
  operating system. For this, use the new `-i` option (or `--install-chef`).
* Add support for Berkshelf v3 in addition to v2. (Thanks to @arosenhagen for
  the original pull request!)
* Re-add long option names like `--host` and `--json-attributes`. Those were
  removed when porting chef-runner to Go (v0.2.0).
* Re-add ability to install chef-runner via Homebrew.

IMPROVEMENTS:

* Add Quick Start guide to README.
* Only prepend `bundle exec` to Ruby commands if Bundler is actually installed.
  A Gemfile alone is no longer enough.
* Add basic [Godoc documentation] to all parts of the source code. As a result,
  [golint] no longer reports any coding style issues.
* Run golint style checks on Travis CI.

BREAKING CHANGES:

* Change default output format of Chef from `null` to `doc`. The former is the
  default used by Vagrant, while the latter is the actual default of Chef Solo.

[Godoc documentation]: https://godoc.org/github.com/mlafeldt/chef-runner
[golint]: https://github.com/golang/lint

## v0.5.0 (Aug 14 2014)

This release brings a couple of improvements to how cookbook dependencies are
resolved, making it possible to use chef-runner as a "general purpose Chef
provisioner" outside of a cookbook directory. (Thanks to @guilhem for the idea!)

For this, chef-runner now follows these rules:

* If the current directory contains a `Berksfile`, Berkshelf is used to manage
  cookbooks.
* If the current directory contains a `Cheffile`, Librarian-Chef is used to
  manage cookbooks.
* If the current directory is a cookbook with a `metadata.rb` file that defines
  the cookbook's name, only this cookbook is copied to the right place.
* The same is done when dependencies have already been resolved with Berkshelf
  or Librarian-Chef to reduce overall provisioning time.
* `metadata.rb` is only required when defining recipes in a format other than
  `cookbook::recipe`. The standard syntax `cookbook::recipe` should always work.
* chef-runner fails if no cookbooks can be found.

Other improvements:

* Configure Chef to verify all HTTPS connections. This removes the SSL warning
  message printed at the start of every Chef run. (Internally, `ssl_verify_mode`
  is set to `:verify_peer`, which is going to be the new default someday.)
* Add `script/coverage` to generate code coverage statistics for Go packages.
* Test dependency handling and provisioning more thoughtfully.
* Log most executed commands when log level is set to debug.
* Improve many parts of the README.
* Add Go documentation to some packages.

## v0.4.0 (Aug 4 2014)

FEATURES:

* Support provisioning of "global" Vagrant machines via their UUID. For this,
  pass the UUID reported by `vagrant global-status` to chef-runner's `-M`
  option. Among other things, this new feature allows you to provision Vagrant
  machines managed by Test Kitchen.
* Use [Librarian-Chef] to install cookbook dependencies if `Cheffile` exists.
  (Also removes temporary Librarian-Chef files from cookbooks before
  transferring them.)
* New option `--version` shows the current program version as well as target OS
  and architecture.

IMPROVEMENTS:

* Extend `script/build` to enable building of download archives with pre-built
  chef-runner binaries for OS X, Linux, FreeBSD, and OpenBSD. Builds triggered
  by this script report the exact Git version that is being compiled.
* Show error message of `vagrant ssh-config` in case it fails.

[Librarian-Chef]: https://github.com/applicationsonline/librarian-chef

## v0.3.0 (Jul 30 2014)

The goal of this release is to *ssh all the things* in order to support any
system reachable over SSH. In addition to local Vagrant machines, chef-runner
can now provision remote machines like EC2 instances. To achieve this, I made
the following changes:

* The argument passed to `-H` now has the format `[user@]hostname[:port]`,
  allowing you to optionally change SSH user and port. (Other SSH settings can
  be set via  `~/.ssh/config`.)
* rsync over SSH is used to transfer files to `/tmp/chef-runner` on the target
  machine. chef-runner no longer depends on `/vagrant` being mounted.
* With Vagrant, instead of running commands via `vagrant ssh`, feed the output
  of `vagrant ssh-config` into OpenSSH. The same SSH configuration is used to
  upload files with rsync.

Other changes:

* Introduce flexible driver concept (inspired by [Test Kitchen]). A driver is
  responsible for running commands on and uploading files to a machine using
  whatever mechanism is available. chef-runner currently contains drivers for
  Vagrant and SSH, but more can -- and will -- be added.
* Always transfer files with `rsync --compress` to speed things up.
* Remove Cucumber scenarios. They didn't add much value to the Go tests and were
  *very* slow. Now Travis builds are much faster.
* More and better Go tests.
* More and better log messages.
* Option `-h` outputs more useful usage text (the one shown in the README).

[Test Kitchen]: https://github.com/test-kitchen/test-kitchen

## v0.2.0 (Jul 18 2014)

This release is a complete rewrite of chef-runner in Go -- a real programming
language that makes it easier to maintain and extend the code base.

The rewrite includes the following changes:

* chef-runner no longer uses Vagrant's Chef configuration. Instead, it creates
  its own local `.chef-runner` folder where configuration data and cookbooks are
  stored. This is the first step towards supporting systems other than Vagrant.
* chef-runner no longer supports long option names like `--host` and
  `--json-attributes`. All original short options are still supported though.
* rsync now only copies actual cookbook files and is run in verbose mode by
  default.
* Events are now properly logged to the console (in color!). The log level can
  be controlled via the `CHEF_RUNNER_LOG` environment variable.
* There are new scripts for bootstrapping, building, and testing the project in
  `script/`.

Hope you enjoy the rewrite. More features to come soon!

## v0.1.2 (Dec 30 2013)

* Add ability to install chef-runner using Homebrew. Thanks to @fh!
* Define environment variable `CHEF_RUNNER_DEBUG` to print commands and their
  arguments as they are executed (`set -x`).
* Fix Cucumber step that handles a command with an option.
* Update documentation.
* Update Cucumber gem.
* Only test against one version of Ruby on Travis.

## v0.1.1 (Nov 29 2013)

* Add Cucumber features.
* Add Travis CI config.
* Move `chef-runner` script to `bin/` folder.
* Document tips on how to further speed things up.

## v0.1.0 (Nov 25 2013)

* Initial public release.
