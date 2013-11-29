Feature: Install cookbook dependencies

  Background:
    Given a cookbook named "cats" with the recipe "default"

  Scenario: Install dependencies with Berkshelf
    When I successfully run `chef-runner`
    Then "berks install --path vendor/cookbooks" should be run
    And /^rsync .*/ should not be run

  Scenario: Update cookbook changes with rsync
    Given a directory named "vendor/cookbooks/cats"
    When I successfully run `chef-runner`
    Then /^rsync .* vendor/cookbooks/cats$/ should be run
    And /^berks .*/ should not be run

  Scenario: Read path to cookbooks from Vagrantfile
    Given a file named "Vagrantfile" with:
    """
      Vagrant.configure("2") do |config|
        config.vm.provision :chef_solo do |chef|
          chef.cookbooks_path = "my/cookbook/folder"
        end
      end
    """
    When I successfully run `chef-runner`
    Then "berks install --path my/cookbook/folder" should be run

  Scenario: Use Bundler when Gemfile is present
    Given an empty file named "Gemfile"
    When I successfully run `chef-runner`
    Then "bundle exec berks install --path vendor/cookbooks" should be run
