class ChefRunner < FPM::Cookery::Recipe
  GOPACKAGE = "github.com/mlafeldt/chef-runner"

  name     "chef-runner"
  version  "0.8.0"
  revision 2
  source   "https://#{GOPACKAGE}/archive/v#{version}.tar.gz"
  sha256   "a7de23f989f8353ecf838b551a8ceff09b83c8aeff2553b2c31d57615f8fcc53"

  description "The fastest way to run Chef cookbooks"
  homepage    "https://#{GOPACKAGE}"
  maintainer  "Mathias Lafeldt <mathias.lafeldt@gmail.com>"
  license     "Apache 2.0"
  section     "development"

  build_depends %w(golang-go git)

  def build
    pkgdir = builddir("gobuild/src/#{GOPACKAGE}")
    mkdir_p pkgdir
    cp_r Dir["*"], pkgdir

    ENV["GOPATH"] = [
      builddir("gobuild/src/#{GOPACKAGE}/Godeps/_workspace"),
      builddir("gobuild"),
    ].join(":")

    safesystem "go version"
    safesystem "go env"
    safesystem "go get -v #{GOPACKAGE}/..."
  end

  def install
    bin.install builddir("gobuild/bin/chef-runner")
  end
end
