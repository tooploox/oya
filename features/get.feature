Feature: Getting packages

Background:
   Given I'm in project dir

Scenario: Get a package
  Given file ./Oyafile containing
    """
    Module: project
    """
  When I run "oya get github.com/bilus/oya@fixtures"
  Then the command succeeds
  And file ./oya/vendor/github.com/bilus/oya/fixtures/features/get.feature/example/Oyafile exists
