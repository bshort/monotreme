package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateTagAbbreviation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple lowercase",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "Simple uppercase",
			input:    "HELLO",
			expected: "hello",
		},
		{
			name:     "Mixed case",
			input:    "HelloWorld",
			expected: "helloworld",
		},
		{
			name:     "Single space",
			input:    "hello world",
			expected: "hello-world",
		},
		{
			name:     "Multiple spaces",
			input:    "hello   world   test",
			expected: "hello-world-test",
		},
		{
			name:     "Special characters",
			input:    "Node.js",
			expected: "nodejs",
		},
		{
			name:     "Multiple special characters",
			input:    "Node.js & JavaScript!",
			expected: "nodejs-javascript",
		},
		{
			name:     "Special chars with spaces",
			input:    "C++ Programming",
			expected: "c-programming",
		},
		{
			name:     "Parentheses",
			input:    "Go (Golang)",
			expected: "go-golang",
		},
		{
			name:     "Underscores",
			input:    "snake_case",
			expected: "snakecase",
		},
		{
			name:     "Hyphens already present",
			input:    "pre-existing-dashes",
			expected: "pre-existing-dashes",
		},
		{
			name:     "Consecutive dashes",
			input:    "test---dash",
			expected: "test-dash",
		},
		{
			name:     "Leading spaces",
			input:    "   leading",
			expected: "leading",
		},
		{
			name:     "Trailing spaces",
			input:    "trailing   ",
			expected: "trailing",
		},
		{
			name:     "Leading and trailing spaces",
			input:    "   both   ",
			expected: "both",
		},
		{
			name:     "Leading special chars",
			input:    "!!!important",
			expected: "important",
		},
		{
			name:     "Trailing special chars",
			input:    "important!!!",
			expected: "important",
		},
		{
			name:     "Numbers",
			input:    "Version 2.0",
			expected: "version-20",
		},
		{
			name:     "Only numbers",
			input:    "123",
			expected: "123",
		},
		{
			name:     "Alphanumeric mix",
			input:    "Test123Test",
			expected: "test123test",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only spaces",
			input:    "     ",
			expected: "",
		},
		{
			name:     "Only special chars",
			input:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "Unicode characters",
			input:    "café résumé",
			expected: "caf-rsum",
		},
		{
			name:     "Complex real-world example 1",
			input:    "Node.js & Express.js Framework",
			expected: "nodejs-expressjs-framework",
		},
		{
			name:     "Complex real-world example 2",
			input:    "C++ Programming (Advanced)",
			expected: "c-programming-advanced",
		},
		{
			name:     "Complex real-world example 3",
			input:    "API Design - RESTful Best Practices",
			expected: "api-design-restful-best-practices",
		},
		{
			name:     "Complex real-world example 4",
			input:    "Python 3.11+ Features",
			expected: "python-311-features",
		},
		{
			name:     "Emoji and special unicode",
			input:    "Test 🚀 Rocket",
			expected: "test-rocket",
		},
		{
			name:     "Forward slashes",
			input:    "OS/2",
			expected: "os2",
		},
		{
			name:     "Backslashes",
			input:    "path\\to\\file",
			expected: "pathtofile",
		},
		{
			name:     "Dots and dashes",
			input:    "file.name-test.txt",
			expected: "filename-testtxt",
		},
		{
			name:     "CamelCase with numbers",
			input:    "HTML5 CSS3 JavaScript",
			expected: "html5-css3-javascript",
		},
		{
			name:     "Quotes and apostrophes",
			input:    "It's a \"test\"",
			expected: "its-a-test",
		},
		{
			name:     "Brackets",
			input:    "[WIP] Feature",
			expected: "wip-feature",
		},
		{
			name:     "Hash and at sign",
			input:    "#hashtag @mention",
			expected: "hashtag-mention",
		},
		{
			name:     "Ampersand with spaces",
			input:    "Rock & Roll",
			expected: "rock-roll",
		},
		{
			name:     "Percent and dollar",
			input:    "100% Success $$$",
			expected: "100-success",
		},
		{
			name:     "Equals and plus",
			input:    "C++ = Great",
			expected: "c-great",
		},
		{
			name:     "Pipes",
			input:    "Option | Another",
			expected: "option-another",
		},
		{
			name:     "Long name with many words",
			input:    "This Is A Very Long Tag Name With Many Words",
			expected: "this-is-a-very-long-tag-name-with-many-words",
		},
		{
			name:     "Single character",
			input:    "A",
			expected: "a",
		},
		{
			name:     "Two characters",
			input:    "AB",
			expected: "ab",
		},
		{
			name:     "Numeric with separators",
			input:    "1,000,000",
			expected: "1000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateTagAbbreviation(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateTagAbbreviation_IdempotentForAlreadyValid(t *testing.T) {
	// Test that already valid abbreviations remain unchanged
	validAbbreviations := []string{
		"hello",
		"hello-world",
		"test123",
		"my-awesome-tag",
		"abc",
		"123",
	}

	for _, abbr := range validAbbreviations {
		t.Run(abbr, func(t *testing.T) {
			result := GenerateTagAbbreviation(abbr)
			require.Equal(t, abbr, result)
		})
	}
}

func TestGenerateTagAbbreviation_ConsistentOutput(t *testing.T) {
	// Test that the same input always produces the same output
	input := "Test Tag Name!"
	expected := GenerateTagAbbreviation(input)

	for i := 0; i < 100; i++ {
		result := GenerateTagAbbreviation(input)
		require.Equal(t, expected, result, "Iteration %d produced different result", i)
	}
}

func TestGenerateTagAbbreviation_SimilarNamesProduceSameAbbreviation(t *testing.T) {
	// Test cases where different inputs should produce the same abbreviation
	testCases := []struct {
		name   string
		inputs []string
	}{
		{
			name: "Different punctuation",
			inputs: []string{
				"Hello World",
				"Hello, World!",
				"Hello - World",
				"Hello. World",
			},
		},
		{
			name: "Different casing",
			inputs: []string{
				"Test Tag",
				"TEST TAG",
				"test tag",
				"TeSt TaG",
			},
		},
		{
			name: "Different spacing",
			inputs: []string{
				"My Tag",
				"My  Tag",
				"My   Tag",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate abbreviation for first input
			expected := GenerateTagAbbreviation(tc.inputs[0])

			// All other inputs should produce the same abbreviation
			for _, input := range tc.inputs[1:] {
				result := GenerateTagAbbreviation(input)
				require.Equal(t, expected, result, "Input '%s' should produce same abbreviation as '%s'", input, tc.inputs[0])
			}
		})
	}
}
