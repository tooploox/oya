Feature: Running hooks

Background:
   Given I'm in project dir

Scenario: Successful run hook
  Given file ./Oyafile containing
    """
    all: |
      foo=4
      if [ $foo -ge 3 ]; then
        touch OK
      fi
      echo "Done"
    """
  When I run "oya run all"
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
  When I run "oya run all"
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

Scenario: No changes
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
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  """

Scenario: Child marks itself as changed
  Given file ./Oyafile containing
    """
    Changeset: echo ""
    all: |
      echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    Changeset: echo "+."
    all: |
      echo "Root"
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Root

  """

Scenario: Child marks parent as changed
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
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Root

  """

Scenario: Parent marks child as changed
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
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Project1

  """

Scenario: No Oyafile
  Given file ./NotOyafile containing
    """
    """
  When I run "oya run all"
  Then the command fails with error
    """
    missing Oyafile
    """

Scenario: Missing hook
  Given file ./Oyafile containing
    """
    """
  When I run "oya run all"
  Then the command fails with error
    """
    missing hook "all"
    """
