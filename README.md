# chef-runner - The fastest way to run Chef cookbooks

[![Build Status](https://travis-ci.org/mlafeldt/chef-runner.svg?branch=master)](https://travis-ci.org/mlafeldt/chef-runner)
[![GoDoc](https://godoc.org/github.com/mlafeldt/chef-runner?status.svg)](https://godoc.org/github.com/mlafeldt/chef-runner)

chef-runner is a tool that speeds up your Chef development and testing workflow.

chef-runner tries hard to provision a machine as fast as possible. It thereby
allows you to change infrastructure code and get *immediate feedback*.

The tool was originally developed as a fast alternative to the painfully slow
`vagrant provision`.

> When it comes to dev/testing, fast feedback is everything. Vagrant needs >5s
> before it even starts provisioning. Each time. Let's fix that. -- [@mlafeldt
> (1 Nov 2013)](https://twitter.com/mlafeldt/status/396299646425137152)

chef-runner has since evolved and can now be used to rapidly provision not only
local Vagrant machines but also remote hosts like EC2 instances.

To further shorten the feedback loop, chef-runner [integrates with
Vim](#use-with-vim) so you don't have to leave your editor while hacking on
recipes.

For more background, check out my blog post *[Telling people about
chef-runner][blog post]*.

## Requirements

To use chef-runner, you need the following software:

* [ssh] command-line tool
* [rsync] command-line tool
* [Chef] must be pre-installed on the target machine

When using chef-runner with [Vagrant], make sure you have a recent version of
Vagrant installed.

The directory you execute chef-runner from must either:

* Include a `Berksfile` so that cookbooks are managed by [Berkshelf]
* Include a `Cheffile` so that cookbooks are managed [Librarian-Chef]
* Be a cookbook with a `metadata.rb` file that defines the cookbook's name

To give you an example, the [Practicing Ruby cookbook][pr-cookbook] is known to
work well with chef-runner.

## Installation

### Download

There are [pre-built binaries] of chef-runner for Mac OS X, Linux, FreeBSD, and
OpenBSD. Please download the proper package for your operating system and
architecture, then unzip the `chef-runner` binary to a location included in
`$PATH`.

### Homebrew

If you're on Mac OS X, the easiest way to get chef-runner is via [Homebrew]:

    $ brew tap mlafeldt/formulas
    $ brew install chef-runner

### Source build

First, make sure you have [Go] version 1.2 or higher.

This single line will download, compile, and install the `chef-runner`
command-line tool:

    $ go get github.com/mlafeldt/chef-runner

For this command to work, `$GOPATH` must be set correctly. Also check that
`$GOPATH` is part of `$PATH`, so that the `chef-runner` executable can be found.
For example, here are the relevant lines from my `~/.bashrc` file:

    export GOPATH="$HOME/devel/go"
    export GOROOT="$(go env GOROOT)"
    export PATH="$GOPATH/bin:$GOROOT/bin:$PATH"

## Usage

### Command Line Reference

chef-runner is a simple command-line tool that has a couple of options:

```
Usage: chef-runner [options] [--] [<recipe>...]

    -h              Show help text
    --version       Show program version

    -H <name>       Set hostname for direct SSH access
    -M <name>       Set name/UUID of Vagrant virtual machine

Options that will be passed to Chef Solo:

    -F <format>     Set output format (null, doc, minimal, min)
                    default: null
    -l <level>      Set log level (debug, info, warn, error, fatal)
                    default: info
    -j <file>       Load attributes from a JSON file
```

### Running Chef Recipes

chef-runner executes one or more recipes you pass on the command line, in the
exact order given. The tool has a flexible recipe syntax allowing you to compose
your run list in multiple ways.

1) Run default recipe when passing no arguments:

    $ chef-runner

2) Run local recipe when passing filename:

    $ chef-runner recipes/foo.rb

3) Run local recipe when passing recipe name:

    $ chef-runner foo

4) Run external recipe when passing `cookbook::recipe`:

    $ chef-runner dogs::bar

5) Run multiple recipes (of mixed type) in order given:

    $ chef-runner recipes/foo.rb bar dogs::baz

**Note: When defining recipes in a format other than `cookbook::recipe`, you
must be inside a cookbook directory with a `metadata.rb` file for chef-runner to
know the cookbook's name.**

Moreover, chef-runner allows you to load node attributes from a local JSON file:

    $ chef-runner -j chef.json

You can also configure both output format and log level of Chef:

    $ chef-runner -F doc -l warn

Last but not least, here is how to enable debug messages:

    $ CHEF_RUNNER_LOG=debug chef-runner ...

### Vagrant

Here's how to use chef-runner with Vagrant machines that are defined in a local
`Vagrantfile` inside the current working directory:

First, make sure that the Vagrant machine you want to provision is running in
the background. You can check the status with `vagrant status`. If the machine
isn't up yet, run `vagrant up`.

If your `Vagrantfile` only defines a single machine, simply run `chef-runner`
and it should work. In a multi-machine environment, use the `-M` option to
specify what Vagrant machine you want to provision. The machine name is the name
you have defined in your `Vagrantfile`. To get a list of all machine names, run
`vagrant status`.

Example:

    $ chef-runner -M db

chef-runner can also provision "global" Vagrant machines that live in a
different directory. For this, all you need to know is the machine's UUID. You
can get a list of all UUIDs by running `vagrant global-status`, e.g.

    $ vagrant global-status
    id       name    provider   state    directory
    -----------------------------------------------------
    a748337  default virtualbox running  /path/to/project
    ...

Then simply pass the UUID of the machine you want to use to chef-runner:

    $ chef-runner -M a748337 ...

Among other things, this allows you to provision Vagrant machines managed by
[Test Kitchen].

### SSH

chef-runner can also provision remote hosts like EC2 instances, or basically
any system reachable over SSH.

Use the `-H` option to specify the name of a host that was configured for direct
SSH access. The argument passed to `-H` has the format `[user@]hostname[:port]`,
allowing you to optionally change SSH user and port. If you need to change other
SSH settings, add a host-specific configuration section to your `~/.ssh/config`.

Examples:

    $ chef-runner -H example.local
    $ chef-runner -H user@example.local
    $ chef-runner -H example.local:1234

### Use with Vim

As a matter of fact, I primarily wrote chef-runner for use with Vim. Instead of
jumping back and forth between editing a Chef recipe and running `vagrant
provision`, I wanted to be able to change some code and get immediate feedback
without having to leave the editor. This is where chef-runner's ability to run a
single recipe file -- the file currently open in Vim -- comes in handy.

There's no Vim plugin (yet). For now, you can just stick this one-liner in your
`.vimrc`:

```vim
nnoremap <leader>r :w\|!chef-runner %<cr>
```

With this key mapping in place, make sure that the target machine is up and open
a recipe in Vim:

```sh
$ vim recipes/default.rb
```

Now whenever you type `<leader>r` (your [leader key] then `r`), chef-runner will
run the *current* recipe on the target machine, giving you fast feedback on
local code changes. (As a bonus, the mapping will also save the file for you.)

Of course, you can also change the key mapping to include whatever chef-runner
options you need. For example, I like using a Chef output format that is less
verbose:

```vim
nnoremap <leader>r :w\|!chef-runner -F min -l warn %<cr>
```

Last but not least, you can always reprogram the Vim key mapping at runtime if
you need to pass project-specific options like host or machine name:

```vim
:nnoremap <leader>r :w\|!chef-runner -F min -l warn -H example.local %<cr>
```

## More Tips

You can further speed up working with chef-runner by doing the following:

* Split up Chef recipes into smaller logical chunks and include those chunks
  using the `include_recipe` method ([good example][pr-recipes]).
* Enable [SSH connection sharing and persistence][ssh-speedup] to speed up
  repeated SSH connections.

## Testing

chef-runner comes with lots of Go unit tests that help to ensure the tool works
as expected. You can run all tests this way:

    $ ./script/test

## License and Author

Author:: Mathias Lafeldt (<mathias.lafeldt@gmail.com>)

Copyright:: 2013-2014, Mathias Lafeldt

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.

## Contributing

Please see `CONTRIBUTING.md` for details.


[Berkshelf]: http://berkshelf.com/
[blog post]: http://mlafeldt.github.io/blog/telling-people-about-chef-runner/
[Chef]: http://www.getchef.com/
[demo video]: http://vimeo.com/78769511
[Go]: http://golang.org/doc/install
[Homebrew]: http://brew.sh/
[leader key]: http://usevim.com/2012/07/20/vim101-leader/
[Librarian-Chef]: https://github.com/applicationsonline/librarian-chef
[pr-cookbook]: https://github.com/elm-city-craftworks/practicing-ruby-cookbook#readme
[pr-recipes]: https://github.com/elm-city-craftworks/practicing-ruby-cookbook/tree/master/recipes
[pre-built binaries]: https://github.com/mlafeldt/chef-runner/releases/latest
[rsync]: http://rsync.samba.org/
[ssh-speedup]: http://interrobeng.com/2013/08/25/speed-up-git-5x-to-50x/
[ssh]: http://www.openssh.com/
[Test Kitchen]: https://github.com/test-kitchen/test-kitchen
[Vagrant]: http://vagrantup.com/
