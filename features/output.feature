Feature: Configure output of Chef

  Background:
    Given a cookbook named "cats" with the recipe "default"

  Scenario: Default output format is "null"
    When I successfully run `chef-runner`
    Then "chef-solo" should be run with the option "--format=null"

  Scenario: Default log level is "info"
    When I successfully run `chef-runner`
    Then "chef-solo" should be run with the option "--log_level=info"

  Scenario: Configure output format via -F
    When I successfully run `chef-runner -F doc`
    Then "chef-solo" should be run with the option "--format=doc"

  Scenario: Configure output format via --format
    When I successfully run `chef-runner --format minimal`
    Then "chef-solo" should be run with the option "--format=minimal"

  Scenario: Configure log level via -l
    When I successfully run `chef-runner -l warn`
    Then "chef-solo" should be run with the option "--log_level=warn"

  Scenario: Configure log level via --log_level
    When I successfully run `chef-runner --log_level debug`
    Then "chef-solo" should be run with the option "--log_level=debug"
