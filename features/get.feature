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
  And file ./.oya/packs/github.com/tooploox/oya-fixtures@v1.0.0/Oyafile exists

Scenario: Get a pack with invalid import
  Given file ./Oyafile containing
    """
    Project: project

    Import:
      invalidPack: foo.com/fooba/fooba
    """
  When I run "oya get github.com/tooploox/oya-fixtures@v1.0.0"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures@v1.0.0/Oyafile exists

Scenario: Get two versions of the same pack
  Given file ./project1/Oyafile containing
    """
    Project: project1

    Require:
      github.com/tooploox/oya-fixtures: v1.0.0

    Import:
      fixtures: github.com/tooploox/oya-fixtures
    """
  And file ./project2/Oyafile containing
    """
    Project: project2

    Require:
      github.com/tooploox/oya-fixtures: v1.1.0

    Import:
      fixtures: github.com/tooploox/oya-fixtures
    """
  When I'm in the ./project1 dir
  And I run "oya run --recurse fixtures.version"
  Then the command succeeds
  And the command outputs
  """
  1.0.0

  """
  When I'm in the ../project2 dir
  And I run "oya run fixtures.version"
  Then the command succeeds
  And the command outputs
  """
  1.1.0

  """

@nocircleci
Scenario: Get a pack from a private repo
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya get github.com/bilus/oya-fixtures-private@v1.0.0"
  Then the command succeeds
  And file ./.oya/packs/github.com/bilus/oya-fixtures-private@v1.0.0/Oyafile exists
