Feature: Getting packs

Background:
   Given I'm in project dir

Scenario: Get a pack
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya get github.com/bilus/oya@fixtures"
  Then the command succeeds
  And file ./oya/vendor/github.com/bilus/oya/fixtures/features/get.feature/example/Oyafile exists
