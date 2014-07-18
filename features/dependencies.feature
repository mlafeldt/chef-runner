Feature: Install cookbook dependencies

  Background:
    Given a cookbook named "cats" with the recipe "default"

  Scenario: Install dependencies with Berkshelf
    When I successfully run `chef-runner`
    Then "berks install --path .chef-runner/cookbooks" should be run
    And /^rsync/ should not be run

  Scenario: Update cookbook changes with rsync
    Given a directory named ".chef-runner/cookbooks/cats"
    When I successfully run `chef-runner`
    Then /^rsync .* .chef-runner/cookbooks/cats$/ should be run
    And /^berks/ should not be run

  Scenario: Use Bundler when Gemfile is present
    Given an empty file named "Gemfile"
    When I successfully run `chef-runner`
    Then "bundle exec berks install --path .chef-runner/cookbooks" should be run
