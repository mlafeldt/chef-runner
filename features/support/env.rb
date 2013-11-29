# Inspired by https://github.com/github/hub/blob/master/features/support/env.rb

require "aruba/cucumber"

bin_dir = File.expand_path("../fakebin", __FILE__)

Before do
  set_env "PATH", "#{bin_dir}:#{ENV['PATH']}"
  set_env "HOME", File.expand_path(File.join(current_dir, "home"))
  FileUtils.mkdir_p ENV["HOME"]
end

World Module.new {
  def history
    histfile = File.join(ENV["HOME"], ".history")
    if File.exist? histfile
      File.readlines(histfile).map(&:chomp)
    else
      []
    end
  end
}
