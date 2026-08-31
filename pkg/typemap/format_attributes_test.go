package typemap_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-seed-data/pkg/typemap"
)

type FormatAttributesTestSuite struct {
	suite.Suite
}

func TestFormatAttributesTestSuite(t *testing.T) {
	suite.Run(t, new(FormatAttributesTestSuite))
}

func (suite *FormatAttributesTestSuite) TestParseFormatHints_Empty() {
	hints := typemap.ParseFormatHints(nil)
	suite.Empty(hints.Format)
	suite.Zero(hints.MinLength)
	suite.Zero(hints.MaxLength)
	suite.Empty(hints.Regex)
	suite.False(hints.HasMin)
	suite.False(hints.HasMax)
}

func (suite *FormatAttributesTestSuite) TestParseFormatHints_Format() {
	hints := typemap.ParseFormatHints([]string{"format=email"})
	suite.Equal("email", hints.Format)
}

func (suite *FormatAttributesTestSuite) TestParseFormatHints_MinMaxLength() {
	hints := typemap.ParseFormatHints([]string{"minLength=5", "maxLength=100"})
	suite.Equal(5, hints.MinLength)
	suite.Equal(100, hints.MaxLength)
}

func (suite *FormatAttributesTestSuite) TestParseFormatHints_Regex() {
	hints := typemap.ParseFormatHints([]string{`regex=^\d{2}-\d{7}$`})
	suite.Equal(`^\d{2}-\d{7}$`, hints.Regex)
}

func (suite *FormatAttributesTestSuite) TestParseFormatHints_MinMax() {
	hints := typemap.ParseFormatHints([]string{"min=1.5", "max=99.9"})
	suite.True(hints.HasMin)
	suite.True(hints.HasMax)
	suite.InDelta(1.5, hints.Min, 0.01)
	suite.InDelta(99.9, hints.Max, 0.01)
}

func (suite *FormatAttributesTestSuite) TestParseFormatHints_IgnoresNonKeyValue() {
	hints := typemap.ParseFormatHints([]string{"optional", "mandatory", "format=phone"})
	suite.Equal("phone", hints.Format)
	suite.Zero(hints.MinLength)
}

func (suite *FormatAttributesTestSuite) TestParseFormatHints_InvalidNumber() {
	hints := typemap.ParseFormatHints([]string{"minLength=abc", "max=xyz"})
	suite.Zero(hints.MinLength)
	suite.False(hints.HasMax)
}

func (suite *FormatAttributesTestSuite) TestHasAttribute() {
	attrs := []string{"optional", "unique", "indexed"}
	suite.True(typemap.HasAttribute(attrs, "optional"))
	suite.True(typemap.HasAttribute(attrs, "unique"))
	suite.False(typemap.HasAttribute(attrs, "mandatory"))
	suite.False(typemap.HasAttribute(nil, "optional"))
}
