package bytesUtils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type test struct {
	line     []byte
	input    []string
	expected boolT
}

func TestContains(t *testing.T) {
	t.Run("Words contains", func(t *testing.T) {
		for _, tst := range [...]test{
			{
				line:     []byte("oleg eglo gelo"),
				input:    []string{"oleg"},
				expected: true,
			},
			{
				line:     []byte("oleg eglo gelo"),
				input:    []string{"noleg"},
				expected: false,
			},
			{
				line:     []byte("127.0.0.1 - - [13/Nov/2020:14:15:07 +0300] \"GET /journal-adm-action/view.Students HTTP/1.1\" 200 9246 \"http://dev.eljur.ru/journal-adm-action\" \"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36\""),
				input:    []string{"Stunts", "Sdents", "Stus", "Students"},
				expected: true,
			},
			{
				line:     []byte(""),
				input:    []string{"Stunts", "Sdents", "Stus", "Students"},
				expected: false,
			},
			{
				line:     []byte(""),
				input:    []string{""},
				expected: true,
			},
		} {
			result := Contains(tst.line, tst.input)
			require.Equal(t, tst.expected, result)
		}
	})

	t.Run("Words not contains", func(t *testing.T) {
		for _, tst := range [...]test{
			{
				line:     []byte("oleg eglo gelo"),
				input:    []string{"oleg"},
				expected: false,
			},
			{
				line:     []byte("oleg eglo gelo"),
				input:    []string{"test"},
				expected: true,
			},
			{
				line:     []byte("oleg eglo gelo"),
				input:    []string{""},
				expected: false,
			},
			{
				line:     []byte(""),
				input:    []string{"oleg"},
				expected: true,
			},
			{
				line:     []byte(""),
				input:    []string{""},
				expected: false,
			},
		} {
			result := Contains(tst.line, tst.input)
			require.Equal(t, !tst.expected, result)
		}
	})
}
