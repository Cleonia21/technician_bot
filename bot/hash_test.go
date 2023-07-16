package bot

import (
	"testing"
)

func TestBot_textToHashToText(t *testing.T) {
	tests := []struct {
		originalText string
	}{
		{
			originalText: "jsdhfgkjfdhk",
		},
		{
			originalText: `dfkh78eyg4i5gy87vjrfvgdf\gf'dgh\fhfg\hdf`,
		},
		{
			originalText: "",
		},
		{
			originalText: "kjehgkjefvlkjefvkljelkrbv\njgfkrewhgferhkjg" +
				"hergfg364tr3whufgy9837fy3iufy93847yfg9843yfg83f" +
				"hergfg364tr3whufgy9837fy3iufy93847yfg9843yfg83f" +
				"hergfg364tr3whufgy9837fy3iufy93847yfg9843yfg83f" +
				"hergfg364tr3whufgy9837fy3iufy93847yfg9843yfg83f" +
				"hergfg364tr3whufgy9837fy3iufy93847yfg9843yfg83f" +
				"hergfg364tr3whufgy9837fy3iufy93847yfg9843yfg83f" +
				"hergfg364tr3whufgy9837fy3iufy93847yfg9843yfg83f" +
				"hergfg364tr3whufgy9837fy3iufy93847yfg9843yfg83f",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			b := &Bot{
				hashCache: make(map[string]string),
			}
			gotHashText := b.textToHash(tt.originalText)
			gotTextFromHash := b.hashToText(gotHashText)
			if tt.originalText != gotTextFromHash {
				t.Errorf("originalText = %v, gotTextFromHash = %v", tt.originalText, gotTextFromHash)
			}
			if len(gotHashText) != 32 {
				t.Errorf("gotHashText len = %v, want = 32", len(gotHashText))
			}
		})
	}
}
