package gherkin

import "github.com/cucumber/messages-go/v10"

// Scenario represents the executed scenario
type Scenario = messages.Pickle

// Step represents the executed step
type Step = messages.Pickle_PickleStep

// DocString represents the DocString argument made to a step definition
type DocString = messages.PickleStepArgument_PickleDocString

// Table represents the Table argument made to a step definition
type Table = messages.PickleStepArgument_PickleTable
