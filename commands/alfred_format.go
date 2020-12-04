package commands

import (
	"encoding/json"
	"fmt"
)

type alfredFormat struct {
	Items []AlfredFormatItem `json:"items"`
}

// AlfredFormatItem represents alfred script filter JSON format item.
type AlfredFormatItem struct {
	UID          string `json:"uid"`
	Title        string `json:"title"`
	SubTitle     string `json:"subtitle"`
	Arg          string `json:"arg"`
	Match        string `json:"match"`
	AutoComplete string `json:"autocomplete"`
}

// AlfredFormatOutput generarete alfred script filter json format string from the given items.
func AlfredFormatOutput(items []AlfredFormatItem, noItemMsg string) (string, error) {
	if len(items) == 0 {
		return fmt.Sprintf("{\"items\": [{\"title\": \"%s\", \"uid\": \"dummy\"}] }", noItemMsg), nil
	}
	alfredFormat := &alfredFormat{
		Items: items,
	}

	output, err := json.Marshal(alfredFormat)
	if err != nil {
		return "", err
	}
	return string(output), nil
}
