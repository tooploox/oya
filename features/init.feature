Feature: Initialization

Background:
   Given I'm in project dir

Scenario: Init a project
  When I run "oya init"
  Then the command succeeds
  And file ./Oyafile exists
