Feature: Manage Secrets for oya

Background:
   Given I'm in project dir

Scenario: It loads values from secrets.oya if present
  Given file ./Oyafile containing
    """
    Project: Secrets
    Values:
      foo: bar

    all: |
      echo ${Oya[foo]}
      echo ${Oya[bar]}
    """
  And file ./secrets.oya containing
    """
    bar: banana
    """
  And I run "oya secrets encrypt"
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  bar
  banana

  """

Scenario: Encrypts secrets file
  Given file ./secrets.oya containing
    """
    foo: SECRETPHRASE
    """
  When I run "oya secrets encrypt"
  Then the command succeeds
  And file ./secrets.oya does not contain
    """
    SECRETPHRASE
    """

Scenario: Views secrets file
  Given file ./secrets.oya containing
    """
    foo: SECRETPHRASE
    """
  And I run "oya secrets encrypt"
  When I run "oya secrets view"
  Then the command succeeds
  And the command outputs to stdout
  """
  foo: SECRETPHRASE
  """

Scenario: It correctly merges secrets
  Given file ./Oyafile containing
    """
    Project: Secrets
    Values:
      foo:
        bar: xxx
        baz: apple

    all: |
      echo ${Oya[foo.bar]}
      echo ${Oya[foo.baz]}
      echo ${Oya[foo.qux]}
    """
  And file ./secrets.oya containing
    """
    foo:
      bar: banana
      qux: peach
    """
  And I run "oya secrets encrypt"
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  banana
  apple
  peach

  """
