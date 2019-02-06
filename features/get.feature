Feature: Getting packs

Background:
   Given I'm in project dir

Scenario: Get a pack
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya get github.com/tooploox/oya-fixtures@v1.0.0"
  Then the command succeeds
  And file ./.oya/vendor/github.com/tooploox/oya-fixtures/Oyafile exists

Scenario: Get a pack with invalid import
  Given file ./Oyafile containing
    """
    Project: project

    Import:
      invalidPack: foo.com/fooba/fooba
    """
  When I run "oya get github.com/tooploox/oya-fixtures@v1.0.0"
  Then the command succeeds
  And file ./.oya/vendor/github.com/tooploox/oya-fixtures/Oyafile exists
