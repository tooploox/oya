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

    render: |
      set -e
      oya render ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
    """
  When I run "oya run foo.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  orange
  """

Scenario: Render uses the scope of the imported Oyafile even with nested imports
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
    Project: pack-foo

    Require:
      github.com/test/bar: v0.0.1

    Import:
      bar: github.com/test/bar

    Values:
      fruit: orange
    """
  And file ./.oya/packs/github.com/test/bar@v0.0.1/Oyafile containing
    """
    Project: pack-bar

    Values:
      fruit: pear

    render: |
      set -e
      oya render ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
    """
  When I run "oya run foo.bar.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  pear
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

    render: |
      set -e
      oya render --auto-scope=false ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
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
      oya run foo.foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Project: pack-foo

    Require:
      github.com/test/bar: v0.0.1

    Import:
      bar: github.com/test/bar

    foo: |
      oya run bar.foo
    """
  And file ./.oya/packs/github.com/test/bar@v0.0.1/Oyafile containing
    """
    Project: pack-bar

    foo: |
      oya run bar

    bar: |
      echo "Success"
    """
  When I run "oya run foo"
  Then the command succeeds
  And the command outputs to stdout
  """
  Success

  """

Scenario: Render called via nested oya runs uses the nested import's scope
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    Values:
      fruit: apple

    foo: |
      oya run foo.foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Project: pack-foo

    Require:
      github.com/test/bar: v0.0.1

    Import:
      bar: github.com/test/bar

    Values:
      fruit: orange

    foo: |
      oya run bar.foo
    """
  And file ./.oya/packs/github.com/test/bar@v0.0.1/Oyafile containing
    """
    Project: pack-bar

    Values:
      fruit: peach

    foo: |
      oya run bar

    bar: |
      set -e
      oya render ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
    """
  When I run "oya run foo"
  Then the command succeeds
  And file ./file.txt contains
  """
  peach
  """

# BUG(bilus): Scenario: Tasks can optionally use the importing Oyafile scope to render values
# BUG(bilus): Scenario: Tasks can optionally use the importing Oyafile scope to run other tasks
