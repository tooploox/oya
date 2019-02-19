Feature: Manage Secrets for oya

Background:
   Given I'm in project dir
   And file ./.sops.yaml containing
    """
    creation_rules:
      - path_regex: ^.*$
        pgp: '317D 6971 DD80 4501 A6B8  65B9 0F1F D46E 2E8C 7202'

    """

Scenario: It loads Values and Tasks from Oyafile.secrets if present
  Given file ./Oyafile containing
    """
    Project: Secrets
    Values:
      foo: bar

    all: |
      echo $foo
      echo $bar
    """
  And file ./secrets.oya containing
    """
    Secrets:
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

# Scenario: Encrypts secrets file
#   Given file ./Oyafile.secrets containing
#     """
#     Values:
#       foo: bar
#     """
#   When I run "oya secrets encrypt"
#   Then the command succeeds
#   And the command outputs to stdout
#   """
#   Done

#   """
#   And file ./OK exists
