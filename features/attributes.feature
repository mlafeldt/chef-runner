Feature: Configure JSON attribute file

  Background:
    Given a cookbook named "cats" with the recipe "default"

  Scenario: Use Vagrant's attribute file by default
    When I successfully run `chef-runner`
    Then "chef-solo" should be run with the option "--json-attributes=/tmp/vagrant-chef-1/dna.json"

  Scenario: Configure attribute file via -j
    Given an empty file named "chef.json"
    When I successfully run `chef-runner -j chef.json`
    Then "chef-solo" should be run with the option "--json-attributes=/vagrant/chef.json"
