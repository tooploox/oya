Feature: Listing available tasks

Background:
   Given I'm in project dir

@current
Scenario: Single Oyafile
  Given file ./Oyafile containing
    """
    Project: project
    build: |
      echo "Done"
    """
  When I run "oya tasks"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya run build

  """

@current
Scenario: Show only user-defined
  Given file ./Oyafile containing
    """
    Project: project
    Changeset: echo +.
    build: |
      echo "Done"
    """
  When I run "oya tasks"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya run build

  """

@current
Scenario: SUbdirectories
  Given file ./Oyafile containing
    """
    Project: project
    build: |
      echo "Done"
    """
  And file ./subdir1/Oyafile containing
    """
    build: |
      echo "Done"
    """
  When I run "oya tasks"
  Then the command succeeds
  And the command outputs to stdout
  """
  # in ./Oyafile
  oya run build

  # in ./subdir1/Oyafile
  oya run build

  """

# @current
# Scenario: Docstring
#   Given file ./Oyafile containing
#     """
#     Project: project
#     build.Doc: Build it
#     build: |
#       echo "Done"
#     """
#   When I run "oya tasks"
#   Then the command succeeds
#   And the command outputs to stdout
#   """
#   # in ./
#   oya run build  # Build it
#   """


# TODO: Subdirs -- execd from project dir
# TODO: Subdirs -- execd from subdir
