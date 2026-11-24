package preserve_nil_map_value

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreserveNilMapValue(t *testing.T) {
	original := &Maps{
		IntItems: map[int32]*Item{
			1: nil,
			2: {},
			3: {Value: 30, Name: "three"},
		},
		StringItems: map[string]*Item{
			"nil": nil,
			"set": {Value: 40, Name: "four"},
		},
		Ints:    map[int32]int32{1: 10},
		Strings: map[string]string{"a": "b"},
	}

	for name, marshal := range map[string]func() ([]byte, error){
		"MarshalVT":       original.MarshalVT,
		"MarshalVTStrict": original.MarshalVTStrict,
	} {
		t.Run(name, func(t *testing.T) {
			size := original.SizeVT()
			data, err := marshal()
			require.NoError(t, err)
			require.Len(t, data, size)

			var decoded Maps
			require.NoError(t, decoded.UnmarshalVT(data))

			require.Len(t, decoded.IntItems, 3)
			require.Nil(t, decoded.IntItems[1])
			require.NotNil(t, decoded.IntItems[2])
			require.Zero(t, decoded.IntItems[2].Value)
			require.NotNil(t, decoded.IntItems[3])
			require.Equal(t, int64(30), decoded.IntItems[3].Value)
			require.Equal(t, "three", decoded.IntItems[3].Name)

			require.Len(t, decoded.StringItems, 2)
			require.Nil(t, decoded.StringItems["nil"])
			require.NotNil(t, decoded.StringItems["set"])
			require.Equal(t, int64(40), decoded.StringItems["set"].Value)
			require.Equal(t, "four", decoded.StringItems["set"].Name)

			require.Equal(t, map[int32]int32{1: 10}, decoded.Ints)
			require.Equal(t, map[string]string{"a": "b"}, decoded.Strings)
		})
	}
}

func TestPreserveNilMapValueNilAndEmptyMaps(t *testing.T) {
	for _, original := range []*Maps{
		{},
		{
			IntItems:    map[int32]*Item{},
			StringItems: map[string]*Item{},
		},
	} {
		size := original.SizeVT()
		data, err := original.MarshalVT()
		require.NoError(t, err)
		require.Len(t, data, size)

		var decoded Maps
		require.NoError(t, decoded.UnmarshalVT(data))
		require.Empty(t, decoded.IntItems)
		require.Empty(t, decoded.StringItems)
		require.Empty(t, decoded.Ints)
		require.Empty(t, decoded.Strings)
	}
}
