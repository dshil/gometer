package example

import (
	"fmt"

	"github.com/dshil/gometer"
	"github.com/gobwas/glob"
)

func ExampleMetricsGetJSONGlobPatterns() {
	metrics := gometer.New()

	for k, v := range map[string]int64{
		"abc":  10,
		"abb":  42,
		"adc":  33,
		"aaac": 17,
	} {
		c := new(gometer.DefaultCounter)
		c.Set(v)
		if err := metrics.Register(k, c); err != nil {
			fmt.Println(err)
			return
		}
	}

	for _, tCase := range [...]struct {
		pattern  string
		expected string
	}{
		{
			pattern:  "*",
			expected: `{"abc": 10, "abb": 42, "adc": 33, "aaac":17}`,
		},
		{
			pattern:  "a*",
			expected: `{"abc": 10, "abb": 42, "adc": 33, "aaac":17}`,
		},
		{
			pattern:  "a?c",
			expected: `{"abc": 10, "adc": 33}`,
		},
		{
			pattern:  "a*c",
			expected: `{"abc": 10, "adc": 33, "aaac":17}`,
		},
		{
			pattern:  "*b*",
			expected: `{"abc": 10, "abb": 42}`,
		},
		{
			pattern:  "??[ab]*",
			expected: `{"abb": 42, "aaac":17}`,
		},
	} {
		g := glob.MustCompile(tCase.pattern)
		b, err := metrics.GetJSON(g.Match)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(b))
	}
	// Output:
	// {"aaac":17,"abb":42,"abc":10,"adc":33}
	// {"aaac":17,"abb":42,"abc":10,"adc":33}
	// {"abc":10,"adc":33}
	// {"aaac":17,"abc":10,"adc":33}
	// {"abb":42,"abc":10}
	// {"aaac":17,"abb":42}
}
