Feature: Built-ins

Background:
   Given I'm in project dir

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
  baz

  """

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
  baz

  """

Scenario: Pack can only run its own tasks
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
      bar: github.com/test/bar
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    foo: |
      echo "foo"
    """
  And file ./.oya/vendor/github.com/test/bar/Oyafile containing
    """
    bar: |
      $Tasks.foo()
    """
  When I run "oya run bar.bar"
  Then the command fails with error matching
    """"
    .*variable not found.*
    """"

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

Scenario: Run render
  Given file ./Oyafile containing
    """
    Project: project
    Values:
      foo: bar

    all: |
      $Render("./templates/file.txt")
    """
  And file ./templates/file.txt containing
    """
    $foo
    """
  When I run "oya run all"
  Then the command succeeds
  And file ./file.txt contains
  """
  bar
  """

Scenario: Run render in alias scope
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo

    Values:
      foo.other: banana
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    Values:
      foo: bar

    all: |
      $Render("$BasePath/templates/file.txt")
    """
  And file ./.oya/vendor/github.com/test/foo/templates/file.txt containing
    """
    $foo
    """
  When I run "oya run foo.all"
  Then the command succeeds
  And file ./file.txt contains
  """
  bar
  """
