Feature: Replace directive

Background:
  Given I'm in project dir


Scenario: Use a local require
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/tooploox/oya-fixtures: v1.0.0

    Replace:
      github.com/tooploox/oya-fixtures: /tmp/pack

    Import:
      foo: github.com/tooploox/oya-fixtures
    """
  And file /tmp/pack/Oyafile containing
    """
    Project: pack

    version: echo 1.0.0
    """
  When I run "oya run foo.version"
  Then the command succeeds
  And the command outputs
    """
    1.0.0

    """
  And file ./.oya/packs/github.com/tooploox/oya-fixtures@v1.0.0/Oyafile does not exist

Scenario: With local require oya doesn't attempt to lookup requirements remotely
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/tooploox/does-not-exist: v1.0.0

    Replace:
      github.com/tooploox/does-not-exist: /tmp/pack

    Import:
      foo: github.com/tooploox/does-not-exist
    """
  And file /tmp/pack/Oyafile containing
    """
    Project: pack

    version: echo 1.0.0
    """
  When I run "oya run foo.version"
  Then the command succeeds

Scenario: Error for Replace in non-root Oyafile
  Given file ./Oyafile containing
    """
    Project: project
    """
  And file ./subdir/Oyafile containing
    """
    Replace:
      github.com/tooploox/oya-fixtures/pack3: ./pack3
    """
  When I'm in the subdir dir
  And I run "oya tasks"
  Then the command fails with error matching
    """
    .*unexpected Replace directive.*
    """
