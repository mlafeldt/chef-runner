# chef-runner - The fastest way to run Chef cookbooks

[![Build Status](https://travis-ci.org/mlafeldt/chef-runner.svg?branch=master)](https://travis-ci.org/mlafeldt/chef-runner)
[![GoDoc](https://godoc.org/github.com/mlafeldt/chef-runner?status.svg)](https://godoc.org/github.com/mlafeldt/chef-runner)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/mlafeldt/chef-runner)

The goal of chef-runner is to speed up your Chef development and testing
workflow by allowing you to change infrastructure code and get *immediate
feedback*.

chef-runner was originally developed as a fast alternative to the painfully slow
`vagrant provision`. The tool has since evolved and can now be used to rapidly
provision not only local Vagrant machines but also remote hosts like EC2
instances.

To further shorten the feedback loop, chef-runner [integrates with Vim][vim] so
you don't have to leave your editor while hacking on recipes.

For more background, check out my blog post ["Telling people about
chef-runner"][blog].

## Quick Start

[Install chef-runner][installation] by either downloading a pre-built binary,
using Homebrew, or running `go get`.

Use chef-runner for local cookbook development with Vagrant:

    $ cd my-awesome-cookook/
    $ vagrant up
    $ chef-runner # will run recipes/default.rb inside the Vagrant machine

Compose Chef run list using flexible recipe syntax:

    $ chef-runner recipes/foo.rb
    $ chef-runner ::foo                        # same as above
    $ chef-runner dogs::bar
    $ chef-runner dogs                         # same as dogs::default
    $ chef-runner recipes/foo.rb bar dogs::baz # will run recipes in given order
    $ chef-runner recipe[cats],dogs::bar       # standard Chef syntax

Provision a specific Vagrant machine in a multi-machine environment:

    $ chef-runner -M db ...

Provision any Vagrant machine by specifying the machine's UUID:

    $ chef-runner -M a748337 ...

Use chef-runner for local cookbook development with Test Kitchen:

    $ kitchen converge default-ubuntu-1404
    $ chef-runner -K default-ubuntu-1404 ...

Use chef-runner as a general purpose Chef provisioner for any system reachable
over SSH:

    $ cd directory-with-berksfile/
    $ chef-runner -H user@example.local apt::default dogs::bar

(chef-runner automatically resolves cookbook dependencies using tools like
Berkshelf or Librarian-Chef.)

If required, install a specific version of Chef before provisioning:

    $ chef-runner -i 11.12.8 ...

## More Information

* See the [chef-runner wiki][wiki] for the official documentation.
* For bug reports and feature requests, please [open an issue here][issues].
* Join our [chat room][gitter] on Gitter for discussion.
* New releases are announced on [Twitter] and the [Chef mailing list][list].

## License

Please see [LICENSE](/LICENSE) for licensing details.

## Want to help?

See the [Development] wiki page for details on how to get the source code and
build chef-runner locally.

## Author

chef-runner is being developed by [Mathias Lafeldt][twitter].


[blog]: http://mlafeldt.github.io/blog/telling-people-about-chef-runner/
[development]: https://github.com/mlafeldt/chef-runner/wiki/Development
[gitter]: https://gitter.im/mlafeldt/chef-runner
[installation]: https://github.com/mlafeldt/chef-runner/wiki/Installation
[issues]: https://github.com/mlafeldt/chef-runner/issues
[list]: http://lists.opscode.com/sympa/info/chef
[twitter]: https://twitter.com/mlafeldt
[vim]: https://github.com/mlafeldt/chef-runner/wiki/Vim
[wiki]: https://github.com/mlafeldt/chef-runner/wiki
