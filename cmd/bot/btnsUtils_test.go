package bot

import (
	"github.com/mymmrac/telego"
	"reflect"
	"testing"
)

func Test_keyboardFromParams(t *testing.T) {
	type args struct {
		rawBtns    []telego.InlineKeyboardButton
		jsonParams string
	}
	tests := []struct {
		name    string
		args    args
		want    []telego.InlineKeyboardButton
		wantErr bool
	}{
		{
			name: "",
			args: args{
				rawBtns: []telego.InlineKeyboardButton{
					{
						Text:         "стартер",
						CallbackData: "стартер",
					},
					{
						Text:         "телематика",
						CallbackData: "телематика",
					},
					{
						Text:         "топливо",
						CallbackData: "топливо",
					},
					{
						Text:         "другое",
						CallbackData: "другое",
					},
					{
						Text:         "акб",
						CallbackData: "акб",
					},
				},
				jsonParams: `{
				  "order": ["телематика", "акб", "топливо", "стартер", "другое"]
				}`,
			},
			want: []telego.InlineKeyboardButton{
				{
					Text:         "телематика",
					CallbackData: "телематика",
				},
				{
					Text:         "акб",
					CallbackData: "акб",
				},
				{
					Text:         "топливо",
					CallbackData: "топливо",
				},
				{
					Text:         "стартер",
					CallbackData: "стартер",
				},
				{
					Text:         "другое",
					CallbackData: "другое",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sortBtns(tt.args.rawBtns, tt.args.jsonParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortBtns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortBtns() got = %v, want %v", got, tt.want)
			}
		})
	}
}
