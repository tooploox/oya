// Package raw implements raw access to Oyafiles, allowing for access
// without actually parsing the files.
// This allows for working with incorrect Oyafiles (e.g. with invalid
// imports) as long as they are well-formed YAML files.
package raw
