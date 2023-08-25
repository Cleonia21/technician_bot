package bot

import (
	"encoding/json"
	"errors"
	"github.com/mymmrac/telego"
	"log"
	"sort"
)

func placementBtns(rawBtns []telego.InlineKeyboardButton, jsonParams string) (*telego.InlineKeyboardMarkup, error) {
	var err error
	if jsonParams != "" {
		rawBtns, err = sortBtns(rawBtns, jsonParams)
		if err != nil {
			return nil, err
		}
	} else {
		sort.Slice(rawBtns, func(i, j int) bool {
			return rawBtns[i].Text < rawBtns[j].Text
		})
	}
	markup := btnsOptimalPlacement(rawBtns)

	if jsonParams == "" {
		keyboard := markup.InlineKeyboard
		sort.Slice(keyboard, func(i, j int) bool {
			return len(keyboard[i]) > len(keyboard[j])
		})
	}
	return markup, nil
}

func sortBtns(rawBtns []telego.InlineKeyboardButton, jsonParams string) ([]telego.InlineKeyboardButton, error) {
	parameterConflictErr := errors.New("discrepancy between the number of buttons and parameters")
	if len(rawBtns) < 2 {
		return nil, parameterConflictErr
	}

	type Order struct {
		Order []string `json:"order"`
	}
	var order Order
	err := json.Unmarshal([]byte(jsonParams), &order)
	if err != nil {
		return nil, err
	}
	sliceOrder := order.Order

	mapOrder := make(map[string]int)
	for key, value := range sliceOrder {
		mapOrder[value] = key
	}

	spacedBtns := make([]telego.InlineKeyboardButton, len(mapOrder))
	for _, btn := range rawBtns {
		i, ok := mapOrder[btn.Text]
		if !ok {
			log.Println("discrepancy between the parameters and the list of buttons") //допились возврат ошибки
			return rawBtns, nil
		}
		spacedBtns[i] = btn
	}
	return spacedBtns, nil
}

// 50 символов в строке
// 2 символа между двумя кнопками
func btnsOptimalPlacement(btns []telego.InlineKeyboardButton) *telego.InlineKeyboardMarkup {
	if len(btns) == 0 {
		return nil
	}
	maxTextLen := 50

	keyboard := &telego.InlineKeyboardMarkup{}
	var row []telego.InlineKeyboardButton
	var rowLen int
	for _, btn := range btns {
		if maxTextLen < rowLen+1+len(btn.Text) {
			if rowLen == 0 {
				row = append(row, btn)
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
				row = []telego.InlineKeyboardButton{}
				continue
			} else {
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
				rowLen = 0
				row = []telego.InlineKeyboardButton{}
			}
		}
		rowLen += 1 + len(btn.Text)
		row = append(row, btn)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	return keyboard
}
