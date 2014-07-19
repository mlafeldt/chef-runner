name              "practicingruby"
maintainer        "Mathias Lafeldt"
maintainer_email  "mathias.lafeldt@gmail.com"
license           "Apache 2.0"
description       "Sets up environment for Practicing Ruby Rails app"
long_description  IO.read(File.join(File.dirname(__FILE__), 'README.md'))
version           "1.3.1"
recipe            "practicingruby::default", "Sets up production-like environment for Practicing Ruby Rails app"

supports "ubuntu", ">= 12.04"

depends "apt", ">= 2.4.0"
depends "database"
depends "mailcatcher"
depends "nginx"
depends "postgresql"
depends "sudo"
depends "user"
