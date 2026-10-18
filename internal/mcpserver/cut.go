package mcpserver

import (
	"encoding/json"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"unicode/utf8"
)

// hintReserve is the room left for the `truncation` object inside the limit, so
// adding it to a cut result does not push the result back over.
const hintReserve = 700

// truncationKey is where a cut result says it was cut. Every output schema
// declares it as an optional property.
const truncationKey = "truncation"

// maxListedCuts bounds the `cuts` array of the hint: it explains the cut, it
// must not become the thing that fills the result.
const maxListedCuts = 6

// What the cut weighs and keeps: a string is a candidate from longStringBytes,
// and a shortening keeps between minKeepRatio and maxKeepRatio of it, never
// less than minKeptBytes.
const (
	longStringBytes = 256
	minKeepRatio    = 0.05
	maxKeepRatio    = 0.9
	minKeptBytes    = 128
)

// truncationSchema describes the hint, for the output schemas. It is the
// property and its type only: the same paragraph on each of 90 tools was a
// tenth of the tool list, so what the hint contains is told once, in the
// server's instructions, and every schema points there (#878).
func truncationSchema() map[string]any {
	return map[string]any{"type": "object"}
}

// withTruncation adds the optional `truncation` property to an output schema.
// The contract's own schemas are closed objects, and a result that carries a
// property its schema does not declare would fail a client's validation.
func withTruncation(schema map[string]any) map[string]any {
	out := make(map[string]any, len(schema)+1)
	maps.Copy(out, schema)

	properties := map[string]any{}
	if existing, ok := schema["properties"].(map[string]any); ok {
		maps.Copy(properties, existing)
	}

	properties[truncationKey] = truncationSchema()
	out["properties"] = properties

	return out
}

// cut is one place where a result was shortened.
type cut struct {
	path  string
	kept  int
	total int
	unit  string
}

// listInfo says that a result is a page of a list under key, so a cut can be
// answered with the page arguments that continue it.
type listInfo struct {
	key     string
	page    int
	perPage int
}

// fit shortens data until its JSON, with text, is at most limit bytes, and says
// so in data's `truncation` property. It works on the decoded value and
// shortens whole array elements or whole strings, so what it hands back is
// always valid JSON, and it cuts the largest part first so the cut falls where
// the bytes are (a list of records, then a very long lesson text) rather than
// on everything alike. It returns false when no cut can reach the limit.
func fit(data map[string]any, textLen, limit int, list *listInfo) bool {
	size := func() int {
		raw, err := json.Marshal(data)
		if err != nil {
			return limit + 1
		}

		return len(raw) + textLen
	}
	if size() <= limit {
		return true
	}

	var cuts []cut

	for range 200 {
		current := size()
		if current+hintReserve <= limit {
			break
		}
		// Aim at the limit in one step; the loop corrects any misjudgement from
		// elements of unequal size.
		ratio := float64(limit-hintReserve) / float64(current)

		target := largest(data)
		if target == nil {
			return false
		}

		cuts = shorten(cuts, target, ratio)
	}

	if size()+hintReserve > limit {
		return false
	}

	hint := map[string]any{"truncated": true, "limitBytes": limit}

	listed := cuts
	if len(listed) > maxListedCuts {
		listed = listed[:maxListedCuts]
	}

	described := make([]any, len(listed))
	for i, c := range listed {
		described[i] = map[string]any{"path": c.path, "kept": c.kept, "total": c.total, "unit": c.unit}
	}

	hint["cuts"] = described
	hint["hint"] = "The response was cut to stay within " + strconv.Itoa(limit) + " bytes. " + continuation(cuts, list, hint)
	data[truncationKey] = hint

	return true
}

// continuation words how to read what was cut, and, for a list, adds the page
// arguments that continue it.
func continuation(cuts []cut, list *listInfo, hint map[string]any) string {
	if list != nil {
		for _, c := range cuts {
			if c.path != list.key || c.unit != "items" {
				continue
			}

			perPage := max(c.kept, 1)
			// The first record of the cut page, as an offset into the whole list,
			// and the page of the smaller size that contains it.
			offset := (list.page - 1) * list.perPage
			page := offset/perPage + 1
			hint["suggestedPerPage"] = perPage
			hint["suggestedPage"] = page

			return fmt.Sprintf("Only %d of %d records of this page fit. Call again with perPage=%d and page=%d, then page=%d and on, to read the rest; the first call may repeat a few records already received.",
				c.kept, c.total, perPage, page, page+1)
		}
	}

	for _, c := range cuts {
		if c.unit == "bytes" {
			return "A text longer than a response may be was shortened at " + c.path + "; the full text is on the website."
		}
	}

	return "A list inside the record was shortened; the full record is on the website."
}

// node is a part of the result that can be shortened, with the way to replace it.
type node struct {
	path  string
	size  int
	value any
	set   func(any)
}

// largest finds the part of data that is the most bytes of JSON among the
// arrays that still have more than one element and the strings long enough to be
// worth shortening.
func largest(data map[string]any) *node {
	var (
		best *node
		walk func(path string, v any, set func(any))
	)

	consider := func(n node) {
		if best == nil || n.size > best.size {
			best = &n
		}
	}
	walk = func(path string, v any, set func(any)) {
		switch x := v.(type) {
		case map[string]any:
			for k, child := range x {
				if path == "" && k == truncationKey {
					continue
				}

				childPath := k
				if path != "" {
					childPath = path + "." + k
				}

				walk(childPath, child, func(nv any) { x[k] = nv })
			}
		case []any:
			if len(x) > 1 {
				// A value decoded from JSON always marshals; were it not, there
				// is no size to weigh it by and it is not a candidate.
				if raw, err := json.Marshal(x); err == nil {
					consider(node{path: path, size: len(raw), value: x, set: set})
				}
			}

			for i, child := range x {
				walk(fmt.Sprintf("%s[%d]", path, i), child, func(nv any) { x[i] = nv })
			}
		case string:
			if len(x) > longStringBytes {
				consider(node{path: path, size: len(x), value: x, set: set})
			}
		}
	}
	walk("", data, nil)

	return best
}

// shorten cuts the node by the ratio, keeping at least one element, or at least
// part of a string, and records what it did (merging with an earlier cut of the
// same place, so the hint reports the first size, not the last).
func shorten(cuts []cut, n *node, ratio float64) []cut {
	ratio = min(max(ratio, minKeepRatio), maxKeepRatio)

	switch x := n.value.(type) {
	case []any:
		keep := min(max(int(float64(len(x))*ratio), 1), len(x)-1)
		n.set(x[:keep])

		return record(cuts, cut{path: n.path, kept: keep, total: len(x), unit: "items"})
	case string:
		keep := min(max(int(float64(len(x))*ratio), minKeptBytes), len(x)-1)
		// Never split a rune: the string must stay valid UTF-8 to stay valid JSON.
		for keep > 0 && !utf8.RuneStart(x[keep]) {
			keep--
		}

		n.set(strings.Clone(x[:keep]) + " [cut]")

		return record(cuts, cut{path: n.path, kept: keep, total: len(x), unit: "bytes"})
	}

	return cuts
}

// record keeps the original total of a place cut more than once.
func record(cuts []cut, c cut) []cut {
	for i := range cuts {
		if cuts[i].path == c.path && cuts[i].unit == c.unit {
			cuts[i].kept = c.kept

			return cuts
		}
	}

	return append(cuts, c)
}
