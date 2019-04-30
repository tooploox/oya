Feature: Running tasks

Background:
   Given I'm in project dir

Scenario: Successfully run task
  Given file ./Oyafile containing
    """
    Project: project
    all: |
      foo=4
      if [ $foo -ge 3 ]; then
        touch OK
      fi
      echo "Done"
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs
  """
  Done

  """
  And file ./OK exists

Scenario: Nested Oyafiles are not processed recursively by default
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
  And the command outputs
  """
  Root

  """
  And file ./Root exists
  And file ./project1/Project1 does not exist
  And file ./project2/Project2 does not exist

Scenario: Nested Oyafiles can be processed recursively
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
  When I run "oya run --recurse all"
  Then the command succeeds
  And the command outputs
  """
  Root
  Project1
  Project2

  """
  And file ./Root exists
  And file ./project1/Project1 exists
  And file ./project2/Project2 exists

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
      foo="${Oya[value]}"
      echo $foo
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs
  """
  some value

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
  And the command outputs
  """
  main

  """

Scenario: Ignore errors in projects inside current project
  Given file ./Oyafile containing
    """
    Project: main
    all: echo "main"
    """
  And file ./foo/Oyafile containing
    """
    Project: foo
    Import:
       xxx: does not exist
    all: echo "foo"
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs
  """
  main

  """

@bug
Scenario: Running recursively
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
  When I run "oya run --recurse all"
  Then the command succeeds
  And the command outputs
  """
  Project1

  """
  And file ./Root does not exist
  And file ./project1/Project1 exists
  And file ./project2/Project2 does not exist

Scenario: Running recursively
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
  When I run "oya run --recurse all"
  Then the command succeeds
  And the command outputs
  """
  Project1

  """
  And file ./Root does not exist
  And file ./project1/Project1 exists
  And file ./project2/Project2 does not exist

Scenario: Running in a subdirectory
  Given file ./Oyafile containing
    """
    Project: project
    all: |
      echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    all: |
      echo "Project1"
    """
  And I'm in the ./project1 dir
  When I run "oya run all"
  Then the command succeeds
  And the command outputs
  """
  Project1

  """

Scenario: Allow empty Require, Import: Values
  Given file ./Oyafile containing
    """
    Project: project

    Require:
    Import:
    Values:

    foo: |
      echo "hello from foo"
    """
  When I run "oya run foo"
  Then the command succeeds
  And the command outputs
  """
  hello from foo

  """

Scenario: Task exits with non-zero code
  Given file ./Oyafile containing
    """
    Project: project

    test: |
      exit 27
    """
  When I run "oya run test"
  Then the command fails
  And the command outputs text matching
    """
    Error: exit status 27

      at line 1, column 1

    > 1\\| exit 27

    """
  And the command exit code is 27

Scenario: Command in task exits with non-zero code
  Given file ./Oyafile containing
    """
    Project: project

    test: |
       bash -c 'exit 27'
       echo "hello after exit 27"
    """
  When I run "oya run test"
  Then the command succeeds
  And the command outputs
  """
  hello after exit 27

  """
  And the command exit code is 0


Scenario: Command in task exits with non-zero code when set -e is in effect
  Given file ./Oyafile containing
    """
    Project: project

    test: |
       set -e
       bash -c 'exit 27'
       echo "hello after exit 27"
    """
  When I run "oya run test"
  Then the command fails
  And the command outputs text matching
    """
    Error: exit status 27

      at line 2, column 1

      1\\| set -e
    > 2\\| bash -c 'exit 27'
    """
  And the command exit code is 27
