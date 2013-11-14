Feature: Override Chef runlist

  Scenario: Run default recipe when passing no arguments
    Given a cookbook named "cats" with the recipe "default"
    When I successfully run `chef-runner`
    Then the runlist should be "cats::default"

  Scenario: Run local recipe when passing filename
    Given a cookbook named "cats" with the recipe "foo"
    When I successfully run `chef-runner recipes/foo.rb`
    Then the runlist should be "cats::foo"

  Scenario: Run local recipe when passing recipe name
    Given a cookbook named "cats" with the recipe "foo"
    When I successfully run `chef-runner foo`
    Then the runlist should be "cats::foo"

  Scenario: Run external recipe when passing cookbook::recipe
    Given a cookbook named "cats"
    When I successfully run `chef-runner dogs::bar`
    Then the runlist should be "dogs::bar"

  Scenario: Run multiple recipes (of mixed type) in order given
    Given a cookbook named "cats" with the recipes "foo,bar"
    When I successfully run `chef-runner recipes/foo.rb bar dogs::baz`
    Then the runlist should be "cats::foo,cats::bar,dogs::baz"
