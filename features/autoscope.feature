Feature: Autoscope

Background:
   Given I'm in project dir

Scenario: Render uses the scope of the imported Oyafile by default
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
      render ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    ${Oya[fruit]}
    """
  When I run "oya run foo.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  orange
  """

Scenario: Render can optionally use the scope of the importing Oyafile
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
      render --auto-scope=false ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    ${Oya[fruit]}
    """
  When I run "oya run foo.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  apple
  """

Scenario: Tasks use the scope of the imported Oyafile to render values
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
       echo ${Oya[fruit]}
    """
  When I run "oya run foo.bar"
  Then the command succeeds
  And the command outputs to stdout
  """
  orange

  """

Scenario: Tasks use the scope of the imported Oyafile to run other tasks
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    foo: |
      run foo.foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Project: pack-foo

    Require:
      github.com/test/bar: v0.0.1

    Import:
      bar: github.com/test/bar

    foo: |
      run bar.foo
    """
  And file ./.oya/packs/github.com/test/bar@v0.0.1/Oyafile containing
    """
    Project: pack-bar

    foo: |
      run bar

    bar: |
      echo "Success"
    """
  When I run "oya run foo"
  Then the command succeeds
  And the command outputs to stdout
  """
  Success

  """

# TODO: Scenario: Tasks can optionally use the importing Oyafile scope to render values
# TODO: Scenario: Tasks can optionally use the importing Oyafile scope to run other tasks
