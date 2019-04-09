Feature: Listing available tasks

Background:
   Given I'm in project dir

Scenario: Single Oyafile
  Given file ./Oyafile containing
    """
    Project: project
    build: |
      echo "Done"
    """
  When I run "oya Oya.tasks"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya build  

  """

Scenario: Show only user-defined
  Given file ./Oyafile containing
    """
    Project: project
    Changeset: echo +.
    build: |
      echo "Done"
    """
  When I run "oya Oya.tasks"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya build  

  """

Scenario: Subdirectories are not recursed by default
  Given file ./Oyafile containing
    """
    Project: project
    build: |
      echo "Done"
    """
  And file ./subdir1/Oyafile containing
    """
    build: |
      echo "Done"
    """
  When I run "oya Oya.tasks"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya build  

  """

Scenario: Subdirectories can be recursed
  Given file ./Oyafile containing
    """
    Project: project
    build: |
      echo "Done"
    """
  And file ./subdir1/Oyafile containing
    """
    build: |
      echo "Done"
    """
  When I run "oya Oya.tasks --recurse"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya build  

  # in ./subdir1/Oyafile
  oya build  

  """

Scenario: Docstring prints
  Given file ./Oyafile containing
    """
    Project: project

    build.Doc: Build it
    build: |
      echo "Done"

    """
  When I run "oya Oya.tasks"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya build  # Build it

  """

Scenario: Doc strings are properly aligned
  Given file ./Oyafile containing
    """
    Project: project

    build.Doc: Build it
    build: |
      echo "Done"

    testAll.Doc: Run all tests
    testAll: |
      echo "Done"
    """
  And file ./subdir1/Oyafile containing
    """
    foo.Doc: Do foo
    foo: |
      echo "Done"
    """
  When I run "oya Oya.tasks --recurse"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya build    # Build it
  oya testAll  # Run all tests

  # in ./subdir1/Oyafile
  oya foo  # Do foo

  """

Scenario: Parent dir tasks are not listed
  Given file ./Oyafile containing
    """
    Project: project

    build.Doc: Build it
    build: |
      echo "Done"

    testAll.Doc: Run all tests
    testAll: |
      echo "Done"
    """
  And file ./subdir1/Oyafile containing
    """
    foo.Doc: Do foo
    foo: |
      echo "Done"
    """
  And I'm in the ./subdir1 dir
  When I run "oya Oya.tasks --recurse"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya foo  # Do foo

  """

Scenario: Imported packs tasks are listed
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    test: |
      echo "Done"
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    packTask: |
      echo "this task is in pack"
    """
  When I run "oya Oya.tasks"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya foo.packTask  
  oya test          

  """
