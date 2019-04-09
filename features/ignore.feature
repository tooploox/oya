Feature: .oyaignore

Background:
   Given I'm in project dir

Scenario: Empty .oyaignore
  Given file ./Oyafile containing
    """
    Project: project
    all: echo "main"
    """
  And file ./oya/subdir/Oyafile containing
    """
    all: echo "subdir"
    """
  When I run "oya --recurse all"
  Then the command succeeds
  And the command outputs to stdout
  """
  main
  subdir

  """

Scenario: Ignore file
  Given file ./Oyafile containing
    """
    Project: project
    Ignore:
      - subdir/Oyafile
    all: echo "main"
    """
  And file ./subdir/Oyafile containing
    """
    all: echo "subdir"
    """
  When I run "oya --recurse all"
  Then the command succeeds
  And the command outputs to stdout
  """
  main

  """

Scenario: Wildcard ignore
  Given file ./Oyafile containing
    """
    Project: project
    Ignore:
      - subdir/*
    all: echo "main"
    """
  And file ./subdir/Oyafile containing
    """
    all: echo "subdir"
    """
  And file ./subdir/foo/Oyafile containing
    """
    all: echo "subdir/foo"
    """
  When I run "oya --recurse all"
  Then the command succeeds
  And the command outputs to stdout
  """
  main

  """
