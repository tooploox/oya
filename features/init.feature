Feature: Initialization

Background:
   Given I'm in project dir

@current
Scenario: Init a project
  When I run "oya init"
  Then the command succeeds
  And file ./Oyafile exists
