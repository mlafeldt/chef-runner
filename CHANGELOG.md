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
