package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func PrintJSON(w io.Writer, value any) error {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(encoded))
	return err
}

func PrintValue(w io.Writer, value any) error {
	switch v := value.(type) {
	case nil:
		_, err := fmt.Fprintln(w, "null")
		return err
	case string:
		_, err := fmt.Fprintln(w, v)
		return err
	case bool, float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, json.Number:
		_, err := fmt.Fprintln(w, v)
		return err
	default:
		return PrintJSON(w, value)
	}
}

func NormalizeInvokeResponse(resp any) any {
	obj, ok := resp.(map[string]any)
	if !ok {
		return resp
	}
	if payload, exists := obj["payload"]; exists {
		return payload
	}
	return resp
}

func PrintSearch(w io.Writer, resp any) error {
	items := extractItems(resp)
	if len(items) == 0 {
		return PrintJSON(w, resp)
	}

	for i, item := range items {
		if row, ok := item.(map[string]any); ok {
			name := firstString(row, "action", "name", "id", "service", "serviceName")
			desc := firstString(row, "description", "summary", "title")
			if name == "" {
				name = jsonCompact(item)
			}
			if desc == "" {
				if _, err := fmt.Fprintf(w, "%d. %s\n", i+1, name); err != nil {
					return err
				}
				continue
			}
			if _, err := fmt.Fprintf(w, "%d. %s - %s\n", i+1, name, desc); err != nil {
				return err
			}
			continue
		}

		if _, err := fmt.Fprintf(w, "%d. %s\n", i+1, jsonCompact(item)); err != nil {
			return err
		}
	}

	return nil
}

func extractItems(resp any) []any {
	switch v := resp.(type) {
	case []any:
		return v
	case map[string]any:
		for _, key := range []string{"results", "actions", "items", "data"} {
			if arr, ok := v[key].([]any); ok {
				return arr
			}
		}
	}
	return nil
}

func firstString(obj map[string]any, keys ...string) string {
	for _, key := range keys {
		if raw, ok := obj[key]; ok {
			if s, ok := raw.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

func jsonCompact(v any) string {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(bytes.TrimSpace(buf.Bytes()))
}
