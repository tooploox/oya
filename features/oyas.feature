Feature: Support for .oya files

Background:
   Given I'm in project dir

Scenario: It loads values from *.oya
  Given file ./Oyafile containing
    """
    Project: Secrets
    Values:
      foo: apple

    showValues: |
      echo ${Oya[foo]}
      echo ${Oya[bar]}
      echo ${Oya[baz]}
    """
  And file ./values1.oya containing
    """
    bar: banana
    """
  And file ./values2.oya containing
    """
    baz: orange
    """
  When I run "oya run showValues"
  Then the command succeeds
  And the command outputs
  """
  apple
  banana
  orange

  """

Scenario: It correctly merges values, processing *.oya alphabetically
  Given file ./Oyafile containing
    """
    Project: Secrets
    Values:
      foo:
        bar: xxx
        baz: apple

    showValues: |
      echo ${Oya[foo.bar]}
      echo ${Oya[foo.baz]}
      echo ${Oya[foo.qux]}
    """
  And file ./0_values.oya containing
    """
    foo:
      bar: banana
      qux: peach
    """
  And file ./1_values.oya containing
    """
    foo:
      qux: coconut
    """
  When I run "oya run showValues"
  Then the command succeeds
  And the command outputs
  """
  banana
  apple
  coconut

  """

Scenario: Support for .oya in a pack
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
    # Project: foo

    echo: |
      echo ${Oya[fruit]}

    """
  And file ./.oya/packs/github.com/test/foo@v0.0.1/values.oya containing
    """
    fruit: orange
    """
  When I run "oya run foo.echo"
  Then the command succeeds
  And the command outputs
  """
  orange

  """
