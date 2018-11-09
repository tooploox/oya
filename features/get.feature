Feature: Getting packages

Background:
   Given I'm in project dir

@current
Scenario: Get a package
  When I run "oya get github.com/bilus/oya@fixtures"
  Then the command succeeds
  And file ./oya/vendor/github.com/bilus/oya/fixtures/features/get.feature/example/Oyafile exists
