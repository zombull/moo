package interactive

import (
	"reflect"
	"sort"
	"strconv"

	"github.com/chzyer/readline"
	"github.com/zombull/floating-castle/bug"
	"github.com/zombull/floating-castle/database"
)

func newSetAutocompleter(s database.Set) readline.AutoCompleter {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	items := make([]readline.PrefixCompleterInterface, len(keys))
	for i, k := range keys {
		items[i] = readline.PcItem(k)
	}
	return readline.NewPrefixCompleter(items...)
}

func newMapAutocompleter(m interface{}) readline.AutoCompleter {
	v := reflect.ValueOf(m)
	bug.On(v.Kind() != reflect.Map || v.Type().Key().Kind() != reflect.String, "passed something other than a map[string]")

	keys := make([]string, len(v.MapKeys()))
	for i, k := range v.MapKeys() {
		keys[i] = k.String()
	}
	sort.Strings(keys)

	items := make([]readline.PrefixCompleterInterface, len(keys))
	for i, k := range keys {
		items[i] = readline.PcItem(k)
	}
	return readline.NewPrefixCompleter(items...)
}

func makeMapInts(start, count int) map[string]int {
	m := make(map[string]int)
	for i := 0; i < count; i++ {
		m[strconv.Itoa(start+i)] = start + i
	}
	return m
}

func makeMapCrags(crags []*database.Crag) map[string]interface{} {
	m := make(map[string]interface{})
	for _, c := range crags {
		m[c.Name] = c
	}
	return m
}

func makeMapAreas(areas []*database.Area) map[string]interface{} {
	m := make(map[string]interface{})
	for _, r := range areas {
		m[r.Name] = r
	}
	return m
}

func makeMapRoutes(routes []*database.Route) map[string]interface{} {
	m := make(map[string]interface{})
	for _, r := range routes {
		m[r.Name] = r
	}
	return m
}

func makeMapSetters(setters []*database.Setter) map[string]interface{} {
	m := make(map[string]interface{})
	for _, s := range setters {
		m[s.Name] = s
	}
	return m
}
