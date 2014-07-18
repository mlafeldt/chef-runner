Feature: Run commands via SSH

  Background:
    Given a cookbook named "cats" with the recipe "default"

  Scenario: Run `vagrant ssh` per default
    When I successfully run `chef-runner`
    Then /^vagrant ssh default -c/ should be run

  Scenario: Run `vagrant ssh` with machine name passed via -M
    When I successfully run `chef-runner -M some-machine`
    Then /^vagrant ssh some-machine -c/ should be run

  Scenario: Run `ssh` with hostname passed via -H
    When I successfully run `chef-runner -H some-host`
    Then /^ssh some-host/ should be run

  Scenario: Disallow passing both host and machine name
    When I run `chef-runner -H some-host -M some-machine`
    Then the exit status should be 1
    And the stderr should contain:
      """
      ERROR: -H and -M cannot be used together
      """
