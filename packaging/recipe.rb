class ChefRunner < FPM::Cookery::Recipe
  GOPACKAGE = "github.com/mlafeldt/chef-runner"

  name     "chef-runner"
  version  "0.9.0"
  revision 1
  source   "https://#{GOPACKAGE}/archive/v#{version}.tar.gz"
  sha256   "4f896fa21cab1f94fe1ce678804b2e5d481523b5c74a5695cbfb76eb9f39dc8b"

  description "The fastest way to run Chef cookbooks"
  homepage    "https://#{GOPACKAGE}"
  maintainer  "Mathias Lafeldt <mathias.lafeldt@gmail.com>"
  license     "Apache 2.0"
  section     "development"

  case platform
  when :debian, :ubuntu
    build_depends %w(git golang-go)
    depends       %w(openssh-client rsync)
  when :centos, :redhat
    build_depends %w(git golang)
    depends       %w(openssh-clients rsync)
  end

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
