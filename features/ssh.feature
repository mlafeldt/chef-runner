Feature: Run commands via SSH

  Background:
    Given a cookbook named "cats" with the recipe "default"

  Scenario: Run `ssh -F` per default
    Given a Vagrant machine named "default"
    When I successfully run `chef-runner`
    Then /^ssh -F .vagrant/machines/default/ssh_config .* default .*/ should be run

  Scenario: Run `ssh -F` with machine name passed via -M
    Given a Vagrant machine named "mybox"
    When I successfully run `chef-runner -M mybox`
    Then /^ssh -F .vagrant/machines/mybox/ssh_config .* mybox .*/ should be run

  Scenario: Run `ssh -F` with machine name passed via --machine
    Given a Vagrant machine named "mybox"
    When I successfully run `chef-runner --machine mybox`
    Then /^ssh -F .vagrant/machines/mybox/ssh_config .* mybox .*/ should be run

  Scenario: Run `ssh` with hostname passed via -H
    When I successfully run `chef-runner -H myhost`
    Then /^ssh myhost .*/ should be run

  Scenario: Run `ssh` with hostname passed via --host
    When I successfully run `chef-runner --host myhost`
    Then /^ssh myhost .*/ should be run

  Scenario: Disallow passing both host and machine name
    When I run `chef-runner --host myhost --machine mybox`
    Then the exit status should be 1
    And the stderr should contain:
      """
      error: --host and --machine cannot be used together
      """
