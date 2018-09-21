Feature: Building

Background:
   Given I'm in project dir

# Scenario: No Oyafile
# Scenario: Missing job

Scenario: Single job
  Given file ./Oyafile containing
    """
    jobs:
      all: |
        foo=4
        if [ $foo -ge 3 ]; then
          touch OK
        fi
        echo "Done"
    """
  When I run "oya build all"
  Then the build succeeds
  # And prints
  # """
  # Done
  # """
  And file ./OK contains
    """
    """


# Scenario: Nested Oyafiles
# Changeset excluding certain dirs
# Scenario: Shell specification
