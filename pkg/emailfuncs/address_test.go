package emailfuncs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRecord struct {
	description string

	raw     string
	invalid bool

	tag    *string
	name   *string
	user   *string
	domain *string
}

func some(str string) *string {
	var p *string
	p = &str
	return p
}

func none() *string {
	return nil
}

var (
	invalidRecord = testRecord{
		description: "Invalid email record",

		raw:     "invalid",
		invalid: true,
	}

	invalidEmptyRecord = testRecord{
		description: "Invalid empty record",

		raw:     "",
		invalid: true,
	}

	edgeCaseRecordWithEmptyTag = testRecord{
		description: "Edge case with empty tag",

		raw: "user+@example.com",

		tag:    none(),
		name:   none(),
		user:   some("user"),
		domain: some("example.com"),
	}

	bareRecord = testRecord{
		description: "Bare user email address",

		raw: "user@example.com",

		tag:    none(),
		name:   none(),
		user:   some("user"),
		domain: some("example.com"),
	}
	bareRecordWithTag = testRecord{
		description: "",

		raw: "user+tag@example.com",

		tag:    some("tag"),
		name:   none(),
		user:   some("user"),
		domain: some("example.com"),
	}

	userRecordWithName = testRecord{
		description: "User email address with name",

		raw: "Bobby Newport <user@example.com>",

		tag:    none(),
		name:   some("Bobby Newport"),
		user:   some("user"),
		domain: some("example.com"),
	}
	userRecordWithNameAndTag = testRecord{
		description: "User email address with name and tag",

		raw: "Bobby Newport <user+tag@example.com>",

		tag:    some("tag"),
		name:   some("Bobby Newport"),
		user:   some("user"),
		domain: some("example.com"),
	}
)

func TestParse(t *testing.T) {
	cases := []testRecord{
		invalidRecord,
		invalidEmptyRecord,

		edgeCaseRecordWithEmptyTag,

		bareRecord,
		bareRecordWithTag,

		userRecordWithName,
		userRecordWithNameAndTag,
	}

	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			a, err := Parse(tc.raw)
			if tc.invalid {
				require.Error(t, err)
				return
			}
			require.NotNil(t, a)

			comparisons := []struct {
				name     string
				expected *string
				fn       func() (string, error)
			}{
				{"tag", tc.tag, a.Tag},
				{"user", tc.user, a.User},
				{"name", tc.name, a.Name},
				{"domain", tc.domain, a.Domain},
			}
			for _, cmp := range comparisons {
				t.Run(cmp.name, func(t *testing.T) {
					actual, err := cmp.fn()
					require.NoError(t, err)
					if cmp.expected != nil {
						assert.Equal(t, *cmp.expected, actual)
						return
					}
					assert.Empty(t, actual)
				})
			}
		})
	}
}

func TestAddressTag(t *testing.T) {
}
