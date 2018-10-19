Feature: Building

Background:
   Given I'm in project dir

Scenario: Successful build
  Given file ./Oyafile containing
    """
    all: |
      foo=4
      if [ $foo -ge 3 ]; then
        touch OK
      fi
      echo "Done"
    """
  When I run "oya build all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Done

  """
  And file ./OK exists


Scenario: Nested Oyafiles
  Given file ./Oyafile containing
    """
    all: |
      touch Root
      echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    all: |
      touch Project1
      echo "Project1"
    """
  And file ./project2/Oyafile containing
    """
    all: |
      touch Project2
      echo "Project2"
    """
  When I run "oya build all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Root
  Project1
  Project2

  """
  And file ./Root exists
  And file ./Project1 exists
  And file ./Project2 exists

Scenario: No rebuild
  Given file ./Oyafile containing
    """
    Changeset: echo ""
    all: |
      echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    Changeset: echo ""
    all: |
      echo "Project1"
    """
  When I run "oya build all"
  Then the command succeeds
  And the command outputs to stdout
  """
  """

@current
Scenario: Child forces parent rebuild
  Given file ./Oyafile containing
    """
    Changeset: echo ""
    all: |
      echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    Changeset: echo "+../"
    all: |
      echo "Root"
    """
  When I run "oya build all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Root

  """

@current
Scenario: Parent forces child rebuild
  Given file ./Oyafile containing
    """
    Changeset: echo "+project1/"
    all: |
      echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    Changeset: echo ""
    all: |
      echo "Project1"
    """
  When I run "oya build all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Project1

  """

# NEXT: Exclusion.

# Scenario: No Oyafile
# Scenario: Missing hook
# Scenario: .oyaignore
# Scenario: Shell specification
# Scenario: Disable early termination
# Scenario: Absolute changeset paths trigger error
# Scenorio/test: Changeset path that doesn't exist
# Scenorio/test: Changeset path that has no Oyafile
