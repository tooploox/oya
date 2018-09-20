Feature: Building


Scenario: Single directory
  Given file Oyafile containing
    """
    all: |
      echo "OK" > OK
    """
  And run oya build all
  # Then the command suceeds
  Then file OK contains
    """
    OK

    """
