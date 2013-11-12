chef-runner
===========

chef-runner explained in one tweet:

> I think I just found a super fast way to edit and run a Chef recipe with Vim
> and Vagrant. Hint: it's not `vagrant provision`. #opschef -- [@mlafeldt (16
> Oct 2013)](https://twitter.com/mlafeldt/status/390235844717838336)

Usage
-----

    Usage: chef-runner [options] [--] [<recipe>...]

        -h, --help                   Show help text
        -H, --host <name>            Set hostname for direct SSH access
        -M, --machine <name>         Set name of Vagrant virtual machine

    Options that will be passed to Chef Solo:

        -F, --format <format>        Set output format (null, doc, minimal, min)
                                     default: null
        -l, --log_level <level>      Set log level (debug, info, warn, error, fatal)
                                     default: info
        -j, --json-attributes <file> Load attributes from a JSON file

Flexible recipe syntax:

    $ chef-runner # runs default recipe
    $ chef-runner recipes/cats.rb
    $ chef-runner cats
    $ chef-runner cats recipes/dogs.rb
    $ chef-runner apt::default
    $ chef-runner apt::default cats recipes/dogs.rb


Vim Integration
---------------

Open recipe in Vim:

```sh
$ vim recipes/ruby.rb
```

Create this key mapping:

```vim
:map ,r :w\|!chef-runner %<cr>
```

Now press `,` + `r` and chef-runner will run the recipe currently open in Vim
inside the Vagrant box, giving you fast feedback on recipe changes.
