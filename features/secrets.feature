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
  And I run "oya secrets encrypt secrets.oya"
  When I run "oya run all"
  Then the command succeeds
  And the command outputs
  """
  bar
  banana

  """

Scenario: Encrypts secrets file
  Given file ./secrets.oya containing
    """
    foo: SECRETPHRASE
    """
  When I run "oya secrets encrypt secrets.oya"
  Then the command succeeds
  And file ./secrets.oya does not contain
    """
    SECRETPHRASE
    """

@bug
Scenario: Views secrets file
  Given file ./secrets.oya containing
    """
    foo: SECRETPHRASE
    """
  Then file ./secrets.oya contains
    """
    foo: SECRETPHRASE
    """
  And I run "oya secrets encrypt secrets.oya"
  Then the command succeeds
  When I run "oya secrets view secrets.oya"
  Then the command succeeds
  And the command outputs
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
  And I run "oya secrets encrypt secrets.oya"
  When I run "oya run all"
  Then the command succeeds
  And the command outputs
  """
  banana
  apple
  peach

  """

Scenario: It can quickly generate and import PGP key
  Given file ./Oyafile containing
    """
    Project: Secrets
    all: |
      echo ${Oya[foo.bar]}
      echo ${Oya[foo.baz]}
    """
  And file ./secrets.oya containing
    """
    foo:
      bar: banana
      baz: peach
    """
  # When I run "oya secrets init --name 'John Public' --email 'john@example.com' --description 'Test key'"
  When I run "oya secrets init"
  Then the command succeeds
  And I run "oya secrets encrypt secrets.oya"
  And I run "oya run all"
  Then the command succeeds
  And the command outputs
  """
  banana
  peach

  """
