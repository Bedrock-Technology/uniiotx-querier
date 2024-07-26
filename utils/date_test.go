package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	timeStr1 := "2023-05-05T14:30:00Z"
	t1, err := time.Parse(time.RFC3339, timeStr1)
	require.Nil(t, err)

	date, y, m, d := Date(t1)

	require.Equal(t, 2023, y)
	require.Equal(t, 5, m)
	require.Equal(t, 5, d)
	require.Equal(t, 20230505, date)

	timeStr2 := "2024-10-13T14:30:00Z"
	t2, err := time.Parse(time.RFC3339, timeStr2)
	require.Nil(t, err)

	date, y, m, d = Date(t2)

	require.Equal(t, 2024, y)
	require.Equal(t, 10, m)
	require.Equal(t, 13, d)
	require.Equal(t, 20241013, date)
}
