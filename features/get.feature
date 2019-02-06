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
  And file ./.oya/vendor/github.com/bilus/oya/fixtures/features/get.feature/example/Oyafile exists

Scenario: Get a pack with invalid import
  Given file ./Oyafile containing
    """
    Project: project

    Import:
      invalidPack: foo.com/fooba/fooba
    """
  When I run "oya get github.com/bilus/oya@fixtures"
  Then the command succeeds
  And file ./.oya/vendor/github.com/bilus/oya/fixtures/features/get.feature/example/Oyafile exists
