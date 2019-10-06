package corpus

import (
	"fmt"
	"testing"
)

var tokens = []Token {
	{
		Term:     "world",
		Position: 1,
		DocID:    1,
		File:     "hamlet and friends",
	},
	{
		Term:     "hamlet",
		Position: 2,
		DocID:    1,
		File:     "hamlet and friends",
	},
	{
		Term:     "world",
		Position: 1,
		DocID:    2,
		File:     "friends and hamlet",
	},
	{
		Term:     "hamlet",
		Position: 2,
		DocID:    2,
		File:     "friends and hamlet",
	},
}

var fileTokens = []Token {
	{
		Term:     "sharkskin",
		Position: 1,
		DocID:    1,
		File:     "hamlet and friends",
	},
	{
		Term:     "hamlet",
		Position: 2,
		DocID:    1,
		File:     "hamlet and friends",
	},
	{
		Term:     "world",
		Position: 1,
		DocID:    2,
		File:     "friends and hamlet",
	},
	{
		Term:     "hamlet",
		Position: 2,
		DocID:    2,
		File:     "friends and hamlet",
	},
}

func TestZoneIndex(t *testing.T) {
	zone := NewZoneIndex()
	zone.BuildZonesIndexFromTokens(tokens, fileTokens)
	fmt.Println(zone.ZoneScore("world", "hamlet"))
}