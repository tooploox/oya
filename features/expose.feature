Feature: Exposing imported package tasks so they can be invoked without the alias.

Background:
  Given I'm in project dir


Scenario: Expose tasks
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
      p: /project1

    Expose: p
    """
  And I'm in the ./project2 dir
  When I run "oya run echo"
  Then the command succeeds
  And the command outputs
  """
  project1

  """

Scenario: Never overwrite existing an task when exposing tasks
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
      p: /project1

    Expose: p

    echo: |
      echo "project2"
    """
  And I'm in the ./project2 dir
  When I run "oya run echo"
  Then the command succeeds
  And the command outputs
  """
  project2

  """

Scenario: Show task as exposed when listing tasks
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
      p: /project1

    Expose: p
    """
  And I'm in the ./project2 dir
  When I run "oya tasks"
  Then the command succeeds
  And the command outputs
  """
  # in ./Oyafile
  oya run echo   # (p.echo)
  oya run p.echo

  """

Scenario: Tasks exposed in imported subdirs carry over when imported
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
      p1: /project1

    Expose: p1
    """
  And file ./project3/Oyafile containing
    """
    Import:
      p2: /project2

    Expose: p2
    """
  And I'm in the ./project3 dir
  When I run "oya run echo"
  Then the command succeeds
  And the command outputs
  """
  project1

  """

Scenario: Tasks exposed in imported packs carry over when imported
  Given file ./Oyafile containing
    """
    Project: main

    Require:
      github.com/tooploox/pack3: v1.0.0

    # Use replace so we don't have to host it on github.
    Replace:
      github.com/tooploox/pack3: /tmp/pack3

    Import:
      pack3: github.com/tooploox/pack3

    Expose: pack3
    """
  And file /tmp/pack1/Oyafile containing
    """
    Project: github.com/tooploox/pack1

    echo: |
      echo "pack1"
    """
  And file /tmp/pack2/Oyafile containing
    """
    Project: github.com/tooploox/pack2

    Require:
      github.com/tooploox/pack1: v1.0.0

    # Use replace so we don't have to host it on github.
    Replace:
      github.com/tooploox/pack1: /tmp/pack1

    Import:
      p1: github.com/tooploox/pack1

    Expose: p1
    """
  And file /tmp/pack3/Oyafile containing
    """
    Project: github.com/tooploox/pack3

    Require:
      github.com/tooploox/pack2: v1.0.0

    # Use replace so we don't have to host it on github.
    Replace:
      github.com/tooploox/pack2: /tmp/pack2

    Import:
      p2: github.com/tooploox/pack2

    Expose: p2
    """
  And I'm in the ./ dir
  When I run "oya run echo"
  Then the command succeeds
  And the command outputs
  """
  pack1

  """
