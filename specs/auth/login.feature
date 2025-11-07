Feature: Login
  As an operator I want to authenticate so that I can manage AD/DHCP

  Scenario: Successful login
    Given the system is running
    When I submit valid credentials
    Then I receive a session cookie

  Scenario: Invalid password
    Given the system is running
    When I submit an invalid password
    Then I get HTTP 401
    And the failed-attempt counter is incremented
