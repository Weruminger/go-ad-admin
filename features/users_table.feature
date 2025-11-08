Feature: UserTabellenVerwaltung
  As an admin
  I want to create and lookup AD users
  So that clients can authenticate

  Scenario: Search with table
    Given an LDAP directory with users:
      | uid         | displayName |
      | anna.smith  | Anna Smith  |
      | joanna.roe  | Joanna Roe  |
      | bob         | Bob Doe     |
    When I search for "anna"
    Then I see 2 results