Feature: Built-ins

Background:
   Given I'm in project dir

Scenario: Access Oyafile base directory
  Given file ./Oyafile containing
    """
    Project: project
    """
  And file ./subdir/Oyafile containing
    """
    all: |
      echo ${Oya[BasePath]}
    """
  When I run "oya run --recurse all"
  Then the command succeeds
  And the command outputs text matching
  """
  ^.*subdir

  """

Scenario: Access pack base directory
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    all: |
      echo ${Oya[BasePath]}
    """
  When I run "oya run foo.all"
  Then the command succeeds
  And the command outputs text matching
  """
  ^.*github.com/test/foo@v0.0.1

  """

Scenario: Access Oyafile Project name
  Given file ./Oyafile containing
    """
    Project: project

    all: |
      echo ${Oya[Project]}
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs text matching
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
      echo ${Oya[Project]}
    """
  When I run "oya run --recurse all"
  Then the command succeeds
  And the command outputs text matching
  """
  project

  """

Scenario: Access Oyafile Project name inside pack
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    all: |
      echo ${Oya[Project]}
    """
  When I run "oya run foo.all"
  Then the command succeeds
  And the command outputs text matching
  """
  project

  """

Scenario: Use plush helpers when rendering
  Given file ./Oyafile containing
    """
    Project: project

    Values:
     arr:
       - 1
       - 2
       - 3

    foo: |
      oya render template.txt
    """
  And file ./template.txt containing
    """
    <%= Len("box") %>
    """
  When I run "oya run foo"
  Then the command succeeds
  And file ./template.txt contains
    """
    3
    """

Scenario: Use sprig functions when rendering (http://masterminds.github.io/sprig)
  Given file ./Oyafile containing
    """
    Project: project

    Values:
     arr:
       - 1
       - 2
       - 3

    foo: |
      oya render template.txt
    """
  And file ./template.txt containing
    """
    <%= Upper(Join(", ", arr)) %>
    """
  When I run "oya run foo bar baz qux"
  Then the command succeeds
  And file ./template.txt contains
    """
    1, 2, 3
    """

Scenario: Use imported pack alias in pack's tasks
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    bar: |
      echo "${Oya[Import.Alias]}"
    """
  When I run "oya run foo.bar"
  Then the command succeeds
  And the command outputs
  """
  foo

  """
