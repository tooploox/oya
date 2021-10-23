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

@current
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


# TODO: Show as aliases when listing tasks.
