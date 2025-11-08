Feature: Privacy modes

  Scenario: Switch privacy to low
    Given privacy mode is "low"
    Then privacy mode should be "low"

  Scenario: Seed with prefix
    Given LDAP contains 3 users with prefix "anna" and display base "Anna"
    When I search for "anna"
    Then I see 3 results

  Scenario: Additive table seeding
    Given LDAP contains users:
      | uid        | displayName |
      | anna.smith | Anna Smith  |
      | bob        | Bob Doe     |
    And LDAP contains users:
      | uid        | displayName |
      | joanna.roe | Joanna Roe  |
    When I search for "anna"
    Then I see 2 results