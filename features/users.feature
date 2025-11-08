Feature: Userverwaltung
  As an admin
  I want to create and lookup AD users
  So that clients can authenticate

  Scenario: User erfolgreich anlegen
    Given an empty directory
    When I create user "alice" with displayName "Alice W."
    Then the user "alice" must exist

  Scenario: Fehler bei fehlendem Display Name
    When I create user "bob" with displayName ""
    Then I receive an error code "INVALID_INPUT"