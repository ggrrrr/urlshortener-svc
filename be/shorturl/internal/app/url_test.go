package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name    string
		fromUrl string
		result  string
		err     error
	}{
		// {
		// 	name:    "ok",
		// 	fromUrl: "http://dir.bg:8080/asdad/adsasd?asdasd=asd",
		// 	result:  "http://dir.bg:8080/asdad/adsasd?asdasd=asd",
		// },
		// {
		// 	name:    "no schema",
		// 	err:     fmt.Errorf("asd"),
		// 	fromUrl: "dir.bg:8080/asdad/adsasd?asdasd=asd",
		// 	result:  "http://dir.bg:8080/asdad/adsasd?asdasd=asd",
		// },
		{
			name:    "no schema",
			err:     fmt.Errorf("asd"),
			fromUrl: "dir.bg",
			result:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseUrl(tc.fromUrl)
			if tc.err == nil {
				require.NoError(t, err)
				require.Equal(t, tc.result, result)
			} else {
				require.Error(t, err)
				fmt.Printf("%v \n", err)
			}
		})
	}
}
