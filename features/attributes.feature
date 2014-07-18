Feature: Configure JSON attribute file

  Background:
    Given a cookbook named "cats" with the recipe "default"

  Scenario: Set attributes to empty hash by default
    When I successfully run `chef-runner`
    Then the file ".chef-runner/dna.json" should contain:
      """
      {}
      """
    And "chef-solo" should be run with the option "--json-attributes /vagrant/.chef-runner/dna.json"

  Scenario: Configure attribute file via -j
    Given a file named "chef.json" with:
      """
      {
        "postgresql": {
          "password": {
            "postgres": "practicingruby"
          }
        }
      }
      """
    When I successfully run `chef-runner -j chef.json`
    Then the file ".chef-runner/dna.json" should contain:
      """
      {
        "postgresql": {
          "password": {
            "postgres": "practicingruby"
          }
        }
      }
      """
    And "chef-solo" should be run with the option "--json-attributes /vagrant/.chef-runner/dna.json"
