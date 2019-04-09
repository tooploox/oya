Feature: Initialization

Background:
   Given I'm in project dir

Scenario: Init a project
  When I run "oya Oya.init"
  Then the command succeeds
  And file ./Oyafile exists

Scenario: Init a existing project
  When I run "oya Oya.init"
  And I run "oya Oya.init"
  Then the command fails with error matching
  """
  .*already an Oya project.*
  """
