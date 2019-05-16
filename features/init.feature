Feature: Initialization

Background:
   Given I'm in project dir

Scenario: Init a project
  When I run "oya init OyaExample"
  Then the command succeeds
  And file ./Oyafile exists

Scenario: Init a existing project
  When I run "oya init OyaExample"
  And I run "oya init OyaExample2"
  Then the command fails with error matching
  """
  .*already an Oya project.*
  """

Scenario: Init a project name
  When I run "oya init OyaExample"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: OyaExample

    """
