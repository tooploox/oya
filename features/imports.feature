Feature: Running tasks

Background:
   Given I'm in project dir

Scenario: Import tasks from installed packs
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    all: |
      touch OK
      echo "Done"
    """
  When I run "oya foo.all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Done

  """
  And file ./OK exists

Scenario: Import task using pack values
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      foo: xxx
    all: |
      bar="${Oya[foo]}"
      echo $bar
    """
  When I run "oya foo.all"
  Then the command succeeds
  And the command outputs to stdout
  """
  xxx

  """

Scenario: Import task using BasePath
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      foo: xxx
    all: |
      bar=$(basename ${Oya[BasePath]})
      echo $bar
    """
  When I run "oya foo.all"
  Then the command succeeds
  And the command outputs to stdout
  """
  foo@v0.0.1

  """

Scenario: Access pack values
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo
    all: |
      echo ${Oya[foo.bar]}
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      bar: xxx
    """
  When I run "oya all"
  Then the command succeeds
  And the command outputs to stdout
  """
  xxx

  """

Scenario: Access current project values
  Given file ./Oyafile containing
    """
    Project: main
    Values:
      foo: main
    """
  And file ./project1/Oyafile containing
    """
    Values:
      foo: project1
    """
  And file ./project2/Oyafile containing
    """
    Import:
      main: /
      p1: /project1
    Values:
      foo: project2
    all: |
      echo ${Oya[main.foo]}
      echo ${Oya[p1.foo]}
      echo ${Oya[foo]}
    """
  When I run "oya --recurse all"
  Then the command succeeds
  And the command outputs to stdout
  """
  main
  project1
  project2

  """

Scenario: Pack values can be set from project Oyafile prefixed with pack alias
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    Values:
      foo:
        fruit: banana
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    all: |
      echo ${Oya[fruit]}
    """
  When I run "oya foo.all"
  Then the command succeeds
  And the command outputs to stdout
  """
  banana

  """

Scenario: Pack values are overriden in main Oyafile
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    Values:
      foo:
        wege: broccoli
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      fruit: banana
      wege: carrot

    all: |
      echo ${Oya[fruit]}
      echo ${Oya[wege]}
    """
  When I run "oya foo.all"
  Then the command succeeds
  And the command outputs to stdout
  """
  banana
  broccoli

  """

# Regression test for #24
Scenario: Import tasks in a subdir Oyafile
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    """
  And file ./subdir/Oyafile containing
    """
    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    all: |
      echo "all"
    """
  And I'm in the ./subdir dir
  When I run "oya --recurse foo.all"
  Then the command succeeds
  And the command outputs to stdout
  """
  all

  """

Scenario: Import tasks from a subdirectory
  Given file ./Oyafile containing
    """
    Project: main
    """
  And file ./project1/Oyafile containing
    """
    Values:
      foo: project1

    echo: |
      echo "project1"
    """
  And file ./project2/Oyafile containing
    """
    Import:
      project1: /project1
    """
  And I'm in the ./project2 dir
  When I run "oya project1.echo"
  Then the command succeeds
  And the command outputs to stdout
  """
  project1

  """
