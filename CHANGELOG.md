## v0.9.0 (unreleased)

* Add ability to provision the host system with `-L` or `--local`. This way, you
  can use chef-runner locally as a convenient wrapper around Chef. Note: This
  is going to replace `-M` (`--machine`) as the default target soon.
* Partial matching of Test Kitchen instance names. For example, `chef-runner -K
  ubuntu` will provision the instance "default-ubuntu-1404" if that's the first
  instance with the string "ubuntu" in its name. Note: The matching does not
  support regular expressions. (Thanks to @StephenKing for the idea.)
* Move most documentation from README to the [chef-runner wiki][wiki]. The
  README was too long, making it too hard to find relevant information. Now it
  only contains the most essential information.
* Add `--resolver` option to specify the cookbook dependency resolver to use.
  Available resolvers are: `berkshelf`, `librarian`, and `dir`. Without this
  option, the resolver is still selected based on the files in the current
  directory.
* Add `--sudo=false` option to not run commands using `sudo`. This currently
  affects Omnibus installer and Chef itself.
* Embed Omnibus installer instead of downloading it from the Internet. Now the
  script [is part][install.sh] of chef-runner's source code and we have total
  control of what it does.
* You can now install chef-runner as a Debian package on most Ubuntu and Debian
  distributions, and as an RPM package on Centos. See the [wiki
  page][wiki-installation] to learn more.
* Add `Makefile`.

[wiki]: https://github.com/mlafeldt/chef-runner/wiki
[wiki-installation]: https://github.com/mlafeldt/chef-runner/wiki/Installation
[install.sh]: /chef/omnibus/assets/install.sh

## v0.8.0 (Nov 16 2014)

FEATURES:

* Support using chef-runner on **Windows**. New releases will include
  cross-compiled Windows binaries. Requires `ssh.exe` and `rsync.exe` to be
  installed. `ssh.exe` is included in MinGW ([Git Bash]). `rsync.exe` must be
  configured to use destination-default permissions when provisioning Unix-like
  systems: `chef-runner --rsync --no-p --rsync --no-g --rsync --chmod=ugo=rwX`.
* Allow to specify one or more custom OpenSSH options on the command line, e.g.
  `chef-runner --ssh LogLevel=debug --ssh "ProxyCommand ..."`. See
  `ssh_config(5)` for a list of available options and their format. (Thanks to
  @berniedurfee who requested this feature.)
* Allow to specify one or more custom Rsync options on the command line, e.g.
  `chef-runner --rsync --progress`. See `rsync(1)` for a list of available
  options.
* Add `--color=false` option to disable colorized output.

IMPROVEMENTS:

* Support standard Chef syntax for composing the run list: entries may be
  separated by comma, and an entry named `foo` will now expand to `foo::default`
  (see BREAKING CHANGES for more information). As a result, something like
  `chef-runner recipe[cats],dogs::bar` now does what Chef users would expect.
* Install Chef using a smart shell wrapper around Omnibus Installer instead of
  running complicated shell commands over SSH. Move installer logic from Chef
  Solo provisioner to new [omnibus package].
* The option `--version` now also outputs the Go version that was used to
  compile chef-runner.

BREAKING CHANGES:

* Adapt run list syntax to Chef's standard: a run list entry named `foo` will
  now expand to `foo::default`. Before, `foo` was treated as a local recipe and
  expanded to `<cookbook>::foo`. Local recipes now need to be passed as `::foo`
  instead. This change also simplifies run list composition when multiple
  cookbooks are involved, e.g. `chef-runner apt postgresql::client nginx`.
* No longer run Rsync in verbose mode by default. To get back the old output,
  you need to use `--rsync --verbose` now.

[Git Bash]: http://msysgit.github.io/
[omnibus package]: https://godoc.org/github.com/mlafeldt/chef-runner/chef/omnibus

## v0.7.0 (Sep 12 2014)

FEATURES:

* Add ability to provision Test Kitchen instances that are reachable over SSH.
  For this, specify the name of a running instance using `-K` (or `--kitchen`).

IMPROVEMENTS:

* Use run list from JSON file if it contains the `run_list` attribute. Recipes
  passed on the command line will still override this list, and
  `recipes/default.rb` is still the default. (Thanks to @arosenhagen who
  requested this feature.)
* Strip all non-cookbook files after resolving dependencies. This ensures that
  only essential cookbook files are copied to target machines, further saving
  time.
* Report overall duration when chef-runner is done.
* Extend `script/build` to auto-generate chef-runner's [Homebrew formula] when
  building a new release with `--release`.
* Extend `script/coverage` to push test coverage statistics from Travis CI to
  [Coveralls].

[Homebrew formula]: https://github.com/mlafeldt/homebrew-formulas/blob/master/Formula/chef-runner.rb
[Coveralls]: https://coveralls.io/r/mlafeldt/chef-runner

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
