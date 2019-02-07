Feature: Built-ins

Background:
   Given I'm in project dir

# Tests for invoking other tasks via $Tasks, do not actually invoke
# oya but rather output the command line. See oya_test.go (MustSetUp)
# for details.

# See oyafile_test.TestRunningTasks for a real example,
# where oya run is actually invoked.
Scenario: Run other tasks
  Given file ./Oyafile containing
    """
    Project: project

    baz: |
      echo "baz"

    bar: |
      echo "bar"
      $Tasks.baz()
    """
  When I run "oya run bar"
  Then the command succeeds
  And the command outputs to stdout
  """
  bar
  oya run baz

  """

# See oyafile_test.TestRunningTasks for a real example,
# where oya run is actually invoked.
Scenario: Run pack's tasks
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    bar: |
      echo "bar"
      $Tasks.baz()

    baz: |
      echo "baz"
    """
  When I run "oya run foo.bar"
  Then the command succeeds
  And the command outputs to stdout
  """
  bar
  oya run foo.baz

  """

# See oyafile_test.TestRunningTasks for a real example,
# where oya run is actually invoked.
Scenario: Run pack's tasks
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    bar: |
      echo "bar"
      $Tasks.baz()

    baz: |
      echo "baz"
    """
  When I run "oya run foo.bar"
  Then the command succeeds
  And the command outputs to stdout
  """
  bar
  oya run foo.baz

  """

Scenario: Access Oyafile base directory
  Given file ./Oyafile containing
    """
    Project: project
    """
  And file ./subdir/Oyafile containing
    """
    all: |
      echo $BasePath
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout text matching
  """
  ^.*subdir

  """

Scenario: Access pack base directory
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    all: |
      echo $BasePath
    """
  When I run "oya run foo.all"
  Then the command succeeds
  And the command outputs to stdout text matching
  """
  ^.*github.com/test/foo

  """

Scenario: Access Oyafile Project name
  Given file ./Oyafile containing
    """
    Project: project

    all: |
      echo $Project
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout text matching
  """
  project

  """

Scenario: Access Oyafile Project name in nested dir
  Given file ./Oyafile containing
    """
    Project: project
    """
  And file ./subdir/Oyafile containing
    """
    all: |
      echo $Project
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout text matching
  """
  project

  """

Scenario: Access Oyafile Project name inside pack
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    all: |
      echo $Project
    """
  When I run "oya run foo.all"
  Then the command succeeds
  And the command outputs to stdout text matching
  """
  project

  """
