Feature: Manage Secrets for oya

Background:
   Given I'm in project dir

Scenario: It loads Values and Tasks from secrets.oya if present
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
    Secrets:
      bar: banana
    """
  And I run "oya Oya.secrets encrypt"
  When I run "oya all"
  Then the command succeeds
  And the command outputs to stdout
  """
  bar
  banana

  """

Scenario: Encrypts secrets file
  Given file ./secrets.oya containing
    """
    Secrets:
      foo: SECRETPHRASE
    """
  When I run "oya Oya.secrets encrypt"
  Then the command succeeds
  And file ./secrets.oya does not contain
    """
    SECRETPHRASE
    """

Scenario: Views secrets file
  Given file ./Oyafile containing
    """
    Project: Secrets
    """
  And file ./secrets.oya containing
    """
    Secrets:
      foo: SECRETPHRASE
    """
  And I run "oya Oya.secrets encrypt"
  When I run "oya Oya.secrets view"
  Then the command succeeds
  And the command outputs to stdout
  """
  Secrets:
    foo: SECRETPHRASE
  """
