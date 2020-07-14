Feature: Global commands

Background:
  Given file ~/Oyafile containing
  """
  foo: |
      echo "foo"
  """

@current
Scenario: Can use the global task in a dir without Oyafile
  Given I'm in an empty dir
  When I run "oya run foo"
  Then the command succeeds
  And the command outputs text matching
