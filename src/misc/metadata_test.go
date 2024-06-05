package misc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadBankMetadata(t *testing.T) {
	t.Parallel()

	banks, err := LoadBankMetadata("../../etc/converted_banks.txt")
	require.NoError(t, err)
	require.NotEmpty(t, banks)
}
