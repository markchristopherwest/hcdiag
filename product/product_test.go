package product

import (
	"testing"

	"github.com/hashicorp/hcdiag/seeker"
	"github.com/stretchr/testify/assert"
)

func TestFilters(t *testing.T) {
	testTable := []struct {
		desc    string
		product *Product
		expect  []*seeker.Seeker
	}{
		{
			desc: "Handles empty seekers and empty filters",
			product: &Product{
				Seekers: []*seeker.Seeker{},
			},
			expect: []*seeker.Seeker{},
		},
		{
			desc: "Handles empty seekers with non-empty filters",
			product: &Product{
				Seekers:  []*seeker.Seeker{},
				Excludes: []string{"hello"},
			},
			expect: []*seeker.Seeker{},
		},
		{
			desc: "Handles nil filters",
			product: &Product{
				Seekers: []*seeker.Seeker{{Identifier: "still here"}},
			},
			expect: []*seeker.Seeker{{Identifier: "still here"}},
		},
		{
			desc: "Handles nil seekers",
			product: &Product{
				Excludes: []string{"nope"},
			},
			expect: []*seeker.Seeker{},
		},
		{
			desc: "Handles empty filters",
			product: &Product{
				Seekers: []*seeker.Seeker{
					{Identifier: "still here"},
				},
				Excludes: []string{},
				Selects:  []string{},
			},
			expect: []*seeker.Seeker{{Identifier: "still here"}},
		},
		{
			desc: "Applies matching excludes",
			product: &Product{
				Seekers: []*seeker.Seeker{
					{Identifier: "goodbye"},
				},
				Excludes: []string{"goodbye"},
			},
			expect: []*seeker.Seeker{},
		},
		{
			desc: "Does not apply non-matching excludes",
			product: &Product{
				Seekers:  []*seeker.Seeker{{Identifier: "goodbye"}},
				Excludes: []string{"hello"},
			},
			expect: []*seeker.Seeker{{Identifier: "goodbye"}},
		},
		{
			desc: "Applies matching Selects",
			product: &Product{
				Seekers: []*seeker.Seeker{
					{Identifier: "goodbye"},
					{Identifier: "hello"},
				},
				Selects: []string{"hello"},
			},
			expect: []*seeker.Seeker{{Identifier: "hello"}},
		},
		{
			desc: "Ignores excludes when Selects are present, and ignores order",
			product: &Product{
				Seekers: []*seeker.Seeker{
					{Identifier: "select3"},
					{Identifier: "select1"},
					{Identifier: "goodbye"},
					{Identifier: "select2"},
				},
				Excludes: []string{"select2", "select3"},
				Selects:  []string{"select2", "select1", "select3"},
			},
			expect: []*seeker.Seeker{
				{Identifier: "select3"},
				{Identifier: "select1"},
				{Identifier: "select2"},
			},
		},
	}

	for _, tc := range testTable {
		err := tc.product.Filter()
		assert.Nil(t, err)
		assert.NotNil(t, tc.product.Seekers)
		assert.Equal(t, tc.expect, tc.product.Seekers, tc.desc)
	}

}

func TestFilterErrors(t *testing.T) {
	testTable := []struct {
		desc    string
		product *Product
		expect  string
	}{
		{
			desc: "Select returns error when pattern is malformed",
			product: &Product{
				Seekers: []*seeker.Seeker{{Identifier: "ignoreme"}},
				Selects: []string{"mal[formed"},
			},
			expect: "filter error: 'syntax error in pattern'",
		},
		{
			desc: "Exclude returns error when pattern is malformed",
			product: &Product{
				Seekers:  []*seeker.Seeker{{Identifier: "ignoreme"}},
				Excludes: []string{"mal[formed"},
			},
			expect: "filter error: 'syntax error in pattern'",
		},
	}

	for _, tc := range testTable {
		err := tc.product.Filter()
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), tc.expect)
		}
	}
}
