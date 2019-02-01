Feature: report results

Scenario: report bad request
  When I supply bad request code
  Then bad request is reported

Scenario: report timeout or error
  When I supply timeouterror code
  Then timeouterror is reported

Scenario: report successful fetch time
  When I supply valid time 100
  Then fetch time is reported
