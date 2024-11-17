package location_test

import (
	"cake-scraper/pkg/location"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MatcherTestSuite struct {
	suite.Suite
}

func (suite *MatcherTestSuite) TestFindBestMatch() {
	testcases := []struct {
		target string
		want   string
	}{
		// Common cases
		{
			target: "Taiwan",
			want:   "Taiwan",
		},
		{
			target: "Taipei City, Taiwan",
			want:   "Taipei City, Taiwan",
		},
		{
			target: "Zhongzheng District, Taipei City, Taiwan",
			want:   "Zhongzheng District, Taipei City, Taiwan",
		},
		// Additional word or omission
		{
			target: "Taipei, Taiwan",
			want:   "Taipei City, Taiwan",
		},
		{
			target: "Zhongzheng District, Taipei City, Taiwan 713",
			want:   "Zhongzheng District, Taipei City, Taiwan",
		},
		{
			target: "Taichung, North District, Taichung City, Taiwan",
			want:   "North District, Taichung City, Taiwan",
		},
		{
			target: "Taichung, North District, Taichung City, Taiwan 404",
			want:   "North District, Taichung City, Taiwan",
		},
		// Empty string
		{
			target: "",
			want:   "",
		},
		// Other countries
		{
			target: "Tokyo, Japan",
			want:   "Japan",
		},
		{
			target: "United States",
			want:   "United States",
		},
		{
			target: "Hong Kong",
			want:   "Hong Kong",
		},
	}
	for _, tc := range testcases {
		suite.Equal(tc.want, location.FindBestMatch(tc.target))
	}
}

func TestMatcherTestSuite(t *testing.T) {
	suite.Run(t, new(MatcherTestSuite))
}
