Feature: Support for pack development

Background:
  Given I'm in project dir


Scenario: Use a local require
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/tooploox/oya-fixtures: v1.0.0

    Replace:
      github.com/tooploox/oya-fixtures: /pack

    Import:
      foo: github.com/tooploox/oya-fixtures
    """
  And file /pack/Oyafile containing
    """
    Project: pack

    version: echo 1.0.0
    """
  When I run "oya run foo.version"
  Then the command succeeds
  And the command outputs to stdout
    """
    1.0.0

    """
  And file ./.oya/packs/github.com/tooploox/oya-fixtures@v1.0.0/Oyafile does not exist
