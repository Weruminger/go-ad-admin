Feature: Search users
  Scenario: Search with result
    Given an LDAP directory with 2 matching users
    When I search for "anna"
    Then I see 2 results
    And no PII is shown when privacy mode is high
