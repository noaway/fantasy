package iconfig

import (
	"os"
	// "regexp"
)

type IConfig interface {
	AddOption(string, string, string) bool
	// RemoveOption removes a option and value from the configuration.
	// It returns true if the option and value were removed, and false otherwise,
	// including if the section did not exist.
	RemoveOption(section string, option string) bool

	// HasOption checks if the configuration has the given option in the section.
	// It returns false if either the option or section do not exist.
	HasOption(section string, option string) bool

	// Options returns the list of options available in the given section.
	// It returns an error if the section does not exist and an empty list if the
	// section is empty. Options within the default section are also included.
	Options(section string) (options []string, err error)

	// SectionOptions returns only the list of options available in the given section.
	// Unlike Options, SectionOptions doesn't return options in default section.
	// It returns an error if the section doesn't exist.
	SectionOptions(section string) (options []string, err error)

	// AddSection adds a new section to the configuration.
	//
	// If the section is nil then uses the section by default which it's already
	// created.
	//
	// It returns true if the new section was inserted, and false if the section
	// already existed.
	AddSection(section string) bool

	// RemoveSection removes a section from the configuration.
	// It returns true if the section was removed, and false if section did not exist.
	RemoveSection(section string) bool

	// HasSection checks if the configuration has the given section.
	// (The default section always exists.)
	HasSection(section string) bool

	// Sections returns the list of sections in the configuration.
	// (The default section always exists).
	Sections() (sections []string)

	// Substitutes values, calculated by callback, on matching regex
	// computeVar(beforeValue *string, regx *regexp.Regexp, headsz, tailsz int, withVar func(*string) string) (*string, error)

	// Bool has the same behaviour as String but converts the response to bool.
	// See "boolString" for string values converted to bool.
	Bool(section string, option string) (value bool, err error)

	// Float has the same behaviour as String but converts the response to float.
	Float(section string, option string) (value float64, err error)

	// Int has the same behaviour as String but converts the response to int.
	Int(section string, option string) (value int, err error)

	// RawString gets the (raw) string value for the given option in the section.
	// The raw string value is not subjected to unfolding, which was illustrated in
	// the beginning of this documentation.
	//
	// It returns an error if either the section or the option do not exist.
	RawString(section string, option string) (value string, err error)

	// RawStringDefault gets the (raw) string value for the given option from the
	// DEFAULT section.
	//
	// It returns an error if the option does not exist in the DEFAULT section.
	RawStringDefault(option string) (value string, err error)

	// String gets the string value for the given option in the section.
	// If the value needs to be unfolded (see e.g. %(host)s example in the beginning
	// of this documentation), then String does this unfolding automatically, up to
	// _DEPTH_VALUES number of iterations.
	//
	// It returns an error if either the section or the option do not exist, or the
	// unfolding cycled.
	String(section string, option string) (value string, err error)

	// WriteFile saves the configuration representation to a file.
	// The desired file permissions must be passed as in os.Open. The header is a
	// string that is saved as a comment in the first line of the file.
	WriteFile(fname string, perm os.FileMode, header string) error
}
