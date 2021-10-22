package gherkin

import messages "github.com/cucumber/messages-go/v16"

// Scenario represents the executed scenario
type Scenario = messages.Pickle

// Step represents the executed step
type Step = messages.PickleStep

// DocString represents the DocString argument made to a step definition
type DocString = messages.PickleDocString

// Table represents the Table argument made to a step definition
type Table = messages.PickleTable
