Feature: Run commands via SSH

  Background:
    Given a cookbook named "cats" with the recipe "default"

  Scenario: Run `vagrant ssh` per default
    When I successfully run `chef-runner`
    Then /^vagrant ssh default -c .*/ should be run

  Scenario: Run `vagrant ssh` with machine name passed via -M
    When I successfully run `chef-runner -M mybox`
    Then /^vagrant ssh mybox -c .*/ should be run

  Scenario: Run `vagrant ssh` with machine name passed via --machine
    When I successfully run `chef-runner --machine mybox`
    Then /^vagrant ssh mybox -c .*/ should be run

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
