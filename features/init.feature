Feature: Initialization

Background:
   Given I'm in project dir

@current
Scenario: Import hooks from vendored idr
  When I run "oya init"
  Then the command succeeds
  And file ./Oyafile exists
