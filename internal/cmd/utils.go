package cmd

import "fmt"

// StringList struct for aggregating string command line flags
type StringList []string

// String TODO: add description
func (l *StringList) String() string {
	return fmt.Sprintf("%v", *l)
}

// Set TODO: add description
func (l *StringList) Set(value string) error {
	*l = append(*l, value)
	return nil
}
