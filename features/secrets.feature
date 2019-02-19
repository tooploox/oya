Feature: Manage Secrets for oya

Background:
   Given I'm in project dir
   And file ./.sops.yaml containing
    """
    creation_rules:
      - path_regex: ^.*$
        pgp: 'D51D DE4C 439A 2FE8 4E67  2AFF 88FF AA75 1034 645E'

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
