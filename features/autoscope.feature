Feature: Autoscope

Background:
   Given I'm in project dir

Scenario: Scope of the importing Oyafile can be optionally used
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    Values:
      fruit: apple
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      fruit: orange

    render:
      $OyaCmd render --auto-scope=false ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    $fruit
    """
  When I run "oya run foo.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  apple
  """

Scenario: Auto scope set to false will run in pack scope
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    Values:
      fruit: apple
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      fruit: orange

    bar:
       echo $fruit
    """
  When I run "oya run --auto-scope=true foo.bar"
  Then the command succeeds
  And the command outputs to stdout
  """
  orange

  """

Scenario: Auto scope will run in importing Oyafile scope
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    foo: |
      echo "Main foo"
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    foo: |
       echo "pack foo"

    bar: |
      $OyaCmd run foo
    """
  When I run "oya run foo.bar"
  Then the command succeeds
  And the command outputs to stdout
  """
  pack foo

  """

Scenario: Auto scope set to false will run in importing Oyafile scope
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    foo: |
      $OyaCmd run foo.foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Project: pack-foo

    Require:
      github.com/test/bar: v0.0.1

    Import:
      bar: github.com/test/bar

    foo: |
      $OyaCmd run bar.foo
    """
  And file ./.oya/packs/github.com/test/bar@v0.0.1/Oyafile containing
    """
    Project: pack-bar

    foo: |
      $OyaCmd run bar

    bar: |
      echo "Success"
    """
  When I run "oya run foo"
  Then the command succeeds
  And the command outputs to stdout
  """
  Success

  """
