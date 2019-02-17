Feature: Rendering templates

Background:
   Given I'm in project dir

Scenario: Render a template
  Given file ./Oyafile containing
    """
    Project: project
    Values:
      foo: xxx
    """
  Given file ./templates/file.txt containing
    """
    $foo
    """
  When I run "oya render ./templates/file.txt"
  Then the command succeeds
  And file ./file.txt contains
  """
  xxx
  """

Scenario: Render a template explicitly pointing to the Oyafile
  Given file ./Oyafile containing
    """
    Project: project
    Values:
      foo: xxx
    """
  Given file ./templates/file.txt containing
    """
    $foo
    """
  When I run "oya render -f ./Oyafile ./templates/file.txt"
  Then the command succeeds
  And file ./file.txt contains
  """
  xxx
  """

Scenario: Render a template directory
  Given file ./Oyafile containing
    """
    Project: project
    Values:
      foo: xxx
      bar: yyy
    """
  Given file ./templates/file.txt containing
    """
    $foo
    """
  Given file ./templates/subdir/file.txt containing
    """
    $bar
    """
  When I run "oya render ./templates/"
  Then the command succeeds
  And file ./file.txt contains
  """
  xxx
  """
  And file ./subdir/file.txt contains
  """
  yyy
  """

Scenario: Render templated paths
  Given file ./Oyafile containing
    """
    Project: project
    Values:
      foo: xxx
      bar: yyy
    """
  Given file ./templates/${foo}.txt containing
    """
    $foo
    """
  Given file ./templates/$bar/${foo}.txt containing
    """
    $bar
    """
  When I run "oya render ./templates/"
  Then the command succeeds
  And file ./xxx.txt contains
  """
  xxx
  """
  And file ./yyy/xxx.txt contains
  """
  yyy
  """

Scenario: Render templated values in alias scope
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
    Values:
      fruit: orange
    """
  And file ./templates/file.txt containing
    """
    $fruit
    """
  When I run "oya render --scope foo ./templates/file.txt"
  Then the command succeeds
  And file ./file.txt contains
  """
  orange
  """

Scenario: Render templated values in alias scope can be overridden
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    Values:
      foo.fruit: banana
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      fruit: orange
    """
  And file ./templates/file.txt containing
    """
    $fruit
    """
  When I run "oya render --scope foo ./templates/file.txt"
  Then the command succeeds
  And file ./file.txt contains
  """
  banana
  """

Scenario: Imported tasks render using target Oyafile scope
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    Values:
      fruit: apple
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      fruit: orange

    render:
      $OyaCmd render ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    $fruit
    """
  When I run "oya run foo.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  apple
  """

Scenario: Alias scope can we detected in imported tasks
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
    Values:
      fruit: orange

    render:
      $OyaCmd render --auto-scope ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    $fruit
    """
  When I run "oya run foo.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  orange
  """

Scenario: Render templated values in alias scope can be overridden when auto-detecting scope
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    Values:
      foo.fruit: banana
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      fruit: orange

    render:
      $OyaCmd render --auto-scope ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    $fruit
    """
  When I run "oya run foo.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  banana
  """
