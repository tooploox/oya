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
    <%= foo %>
    """
  When I run "oya Oya.render ./templates/file.txt"
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
    <%= foo %>
    """
  When I run "oya Oya.render -f ./Oyafile ./templates/file.txt"
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
    <%= foo %>
    """
  Given file ./templates/subdir/file.txt containing
    """
    <%= bar %>
    """
  When I run "oya Oya.render ./templates/"
  Then the command succeeds
  And file ./file.txt contains
  """
  xxx
  """
  And file ./subdir/file.txt contains
  """
  yyy
  """

Scenario: Rendering values in specified scope pointing to imported pack
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
    <%= fruit %>
    """
  When I run "oya Oya.render --scope foo ./templates/file.txt"
  Then the command succeeds
  And file ./file.txt contains
  """
  orange
  """

Scenario: Rendered values in specified scope can be overridden
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    Values:
      foo:
        fruit: banana
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      fruit: orange
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
    """
  When I run "oya Oya.render --scope foo ./templates/file.txt"
  Then the command succeeds
  And file ./file.txt contains
  """
  banana
  """


Scenario: Imported tasks render using their own Oyafile scope by default
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

    render: |
      oya Oya.render ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
    """
  When I run "oya foo.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  orange
  """

Scenario: Values in imported pack scope can be overridden
  Given file ./Oyafile containing
    """
    Project: project

    Require:
      github.com/test/foo: v0.0.1

    Import:
      foo: github.com/test/foo

    Values:
      foo:
        fruit: banana
    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/Oyafile containing
    """
    Values:
      fruit: orange

    render:
      oya Oya.render ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
    """
  When I run "oya foo.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  banana
  """
Scenario: Scope of the importing Oyafile can be optionally used
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
      oya Oya.render --auto-scope=false ./templates/file.txt
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
    """
  When I run "oya foo.render"
  Then the command succeeds
  And file ./file.txt contains
  """
  apple
  """


Scenario: Rendering values in specified scope
  Given file ./Oyafile containing
    """
    Project: project

    Values:
      fruits:
        fruit: orange
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
    """
  When I run "oya Oya.render --scope fruits ./templates/file.txt"
  Then the command succeeds
  And file ./file.txt contains
  """
  orange
  """

Scenario: Rendering values in specified nested scope
  Given file ./Oyafile containing
    """
    Project: project

    Values:
      plants:
        fruits:
          fruit: orange
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
    """
  When I run "oya Oya.render --scope plants.fruits ./templates/file.txt"
  Then the command succeeds
  And file ./file.txt contains
  """
  orange
  """

Scenario: Rendering one file to an output dir
  Given file ./Oyafile containing
    """
    Project: project

    Values:
      fruit: orange
    """
  And file ./templates/file.txt containing
    """
    <%= fruit %>
    """
  When I run "oya Oya.render --output-dir ./foobar ./templates/file.txt"
  Then the command succeeds
  And file ./foobar/file.txt contains
  """
  orange
  """

Scenario: Rendering a dir to an output dir
  Given file ./Oyafile containing
    """
    Project: project

    Values:
      culprit: Eve
      weapon: apple
    """
  And file ./templates/file1.txt containing
    """
    <%= weapon %>
    """
  And file ./templates/file2.txt containing
    """
    <%= culprit %>
    """
  When I run "oya Oya.render --output-dir ./foobar ./templates/"
  Then the command succeeds
  And file ./foobar/file1.txt contains
  """
  apple
  """
  And file ./foobar/file2.txt contains
  """
  Eve
  """

Scenario: Render dir excluding files and directories
  Given file ./Oyafile containing
    """
    Project: project
    Values:
      foo: xxx
      bar: yyy
    """
  Given file ./templates/file.txt containing
    """
    <%= foo %>
    """
  And file ./templates/excludeme.txt containing
    """
    <%= badvariable %>
    """
  And file ./templates/subdir/excludeme.txt containing
    """
    <%= badvariable %>
    """
  And file ./templates/subdir/file.txt containing
    """
    <%= bar %>
    """
  And file ./templates/excludeme/excludeme.txt containing
    """
    <%= badvariable %>
    """
  When I run "oya Oya.render --exclude excludeme.txt --exclude subdir/excludeme.txt --exclude excludeme/* ./templates/"
  Then the command succeeds
  And file ./file.txt contains
  """
  xxx
  """
  And file ./subdir/file.txt contains
  """
  yyy
  """
  And file ./excludeme.txt does not exist
  And file ./subdir/excludeme.txt does not exist
  And file ./excludeme/excludeme.txt does not exist


Scenario: Render dir excluding using globbing
  Given file ./Oyafile containing
    """
    Project: project
    Values:
      foo: xxx
      bar: yyy
    """
  Given file ./templates/file.txt containing
    """
    <%= foo %>
    """
  And file ./templates/excludeme.txt containing
    """
    <%= badvariable %>
    """
  And file ./templates/subdir/excludeme.txt containing
    """
    <%= badvariable %>
    """
  And file ./templates/subdir/file.txt containing
    """
    <%= bar %>
    """
  And file ./templates/excludeme/excludeme.txt containing
    """
    <%= badvariable %>
    """
  When I run "oya Oya.render --exclude **excludeme.txt ./templates/"
  Then the command succeeds
  And file ./file.txt contains
  """
  xxx
  """
  And file ./subdir/file.txt contains
  """
  yyy
  """
  And file ./excludeme.txt does not exist
  And file ./subdir/excludeme.txt does not exist
  And file ./excludeme/excludeme.txt does not exist


Scenario: Rendering a dir to an output dir outside project dir
  Given file ./Oyafile containing
    """
    Project: project

    Values:
      culprit: Eve
      weapon: apple
    """
  And file ./templates/file1.txt containing
    """
    <%= weapon %>
    """
  And file ./templates/file2.txt containing
    """
    <%= culprit %>
    """
  When I run "oya Oya.render --output-dir /tmp/foobar ./templates/"
  Then the command succeeds
  And file /tmp/foobar/file1.txt contains
  """
  apple
  """
  And file /tmp/foobar/file2.txt contains
  """
  Eve
  """
