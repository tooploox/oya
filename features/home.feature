Feature: Home directory

Background:
   Given I'm in project dir

Scenario: Decide where installed packs are stored
  Given file ./Oyafile containing
    """
    Project: project
    """
  And the OYA_HOME environment variable set to "/tmp/oya_home"
  When I run "oya Oya.get github.com/tooploox/oya-fixtures@v1.0.0"
  Then the command succeeds
  And file /tmp/oya_home/.oya/packs/github.com/tooploox/oya-fixtures@v1.0.0/Oyafile exists
