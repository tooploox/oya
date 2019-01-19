Feature: Running tasks

Background:
   Given I'm in project dir

Scenario: Successfully run task
  Given file ./Oyafile containing
    """
    Project: project
    all: |
      foo=4
      if [ $$foo -ge 3 ]; then
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


@xxx
Scenario: Nested Oyafiles
  Given file ./Oyafile containing
    """
    Project: project
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
  And file ./project1/Project1 exists
  And file ./project2/Project2 exists

Scenario: No changes
  Given file ./Oyafile containing
    """
    Project: project
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
    Project: project
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
    Project: project
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
    Project: project
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
    Project: project
    """
  When I run "oya run all"
  Then the command fails with error matching
  """
  .*no Oyafile project in.*
  """

Scenario: Missing task
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya run all"
  Then the command fails with error
    """
    missing task "all"
    """

Scenario: Script template
  Given file ./Oyafile containing
    """
    Project: project
    Values:
      value: some value
    all: |
      foo="$value"
      echo $$foo
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  some value

  """

Scenario: Ignore vendored Oyafiles
  Given file ./Oyafile containing
    """
    Project: project
    all: echo "main"
    """
  And file ./oya/vendor/github.com/test/foo/Oyafile containing
    """
    all: echo "vendored"
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  main

  """

Scenario: Ignore projects inside current project
  Given file ./Oyafile containing
    """
    Project: main
    all: echo "main"
    """
  And file ./foo/Oyafile containing
    """
    Project: foo
    all: echo "foo"
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  main

  """

@xxx
Scenario: Running in subdir
  Given file ./Oyafile containing
    """
    Project: project
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
  And I'm in the ./project1 dir
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Project1

  """
  And file ././Root does not exist
  And file ./Project1 exists
  And file ./../project2/Project2 does not exist
