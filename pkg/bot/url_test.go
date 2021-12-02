package bot

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHasYoutubeLink(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{input: "hello world", want: nil},
		{input: "http://www.youtube.com/watch?v=iwGFalTRHDA", want: []string{"iwGFalTRHDA"}},
		{input: "http://www.youtube.com/watch?v=iwGFalTRHDA&feature=related", want: []string{"iwGFalTRHDA"}},
		{input: "http://youtu.be/iwGFalTRHDA", want: []string{"iwGFalTRHDA"}},
		{input: "http://youtu.be/n17B_uFF4cA", want: []string{"n17B_uFF4cA"}},
		{input: "http://www.youtube.com/embed/watch?feature=player_embedded&v=r5nB9u4jjy4", want: []string{"r5nB9u4jjy4"}},
		{input: "http://www.youtube.com/watch?v=t-ZRX8984sc", want: []string{"t-ZRX8984sc"}},
		{input: "http://youtu.be/t-ZRX8984sc", want: []string{"t-ZRX8984sc"}},
		{input: "http://youtu.be/t-ZRX8984sc and http://youtu.be/iwGFalTRHDA", want: []string{"t-ZRX8984sc", "iwGFalTRHDA"}},
	}

	for i, tt := range tests {
		q := tt
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, hasYoutubeLink(q.input), q.want)
		})
	}
}
