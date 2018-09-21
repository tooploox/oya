Feature: Building

# Scenario: No Oyafile
# Scenario: Missing job

Scenario: Single job
  Given file Oyafile containing
    """
    jobs:
      all: |
        foo=4
        if [ $foo -ge 3 ]; then
          touch OK
        fi
        echo "Done"
    """
  When "oya build all" is run
  # Then the command suceeds
  # And prints
  # """
  # Done
  # """
  Then file OK contains
    """
    """


# Scenario: Nested Oyafiles

# Scenario: Shell specification
