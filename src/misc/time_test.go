package misc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomTime(t *testing.T) {
	t.Parallel()

	t.Run("parse date time successfully", func(t *testing.T) {
		t.Parallel()

		data := []byte("2023-11-19 09:30:00")
		var unmarshalTime DateTime
		require.NoError(t, unmarshalTime.UnmarshalJSON(data))
	})
}
