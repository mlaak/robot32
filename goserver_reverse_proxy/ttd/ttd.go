//go:build ttd

package ttd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
)

type OrderedMap struct {
	keys   []string
	values map[string]interface{}
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		keys:   make([]string, 0),
		values: make(map[string]interface{}),
	}
}

func (om *OrderedMap) Set(key string, value interface{}) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

func (om *OrderedMap) Get(key string) (interface{}, bool) {
	val, exists := om.values[key]
	return val, exists
}

func (om *OrderedMap) Delete(key string) {
	if _, exists := om.values[key]; exists {
		delete(om.values, key)
		for i, k := range om.keys {
			if k == key {
				om.keys = append(om.keys[:i], om.keys[i+1:]...)
				break
			}
		}
	}
}

func (om *OrderedMap) Keys() []string {
	return om.keys
}

func (om *OrderedMap) Len() int {
	return len(om.keys)
}

func (om *OrderedMap) Range(f func(key string, value interface{}) bool) {
	for _, key := range om.keys {
		if !f(key, om.values[key]) {
			break
		}
	}
}

// MarshalJSON implements the json.Marshaler interface
func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	// Create a slice of map[string]interface{} to preserve order
	var pairs []map[string]interface{}

	for _, key := range om.keys {
		pairs = append(pairs, map[string]interface{}{key: om.values[key]})
	}

	return json.Marshal(pairs)
}

type fileLine struct {
	FN       string
	LN       int
	M        string
	P        *OrderedMap
	Path     string
	Parent   *fileLine
	Filedata string
}

func filePath(s string) string {

	u, err := url.Parse(s)
	if err != nil {
		fmt.Println("Error parsing Path:", err)
		return ""
	}
	return path.Base(u.Path)
}

func appendToFile(filename, text string) error {
	// Open the file in append mode, or create it if it doesn't exist
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the text to the file
	if _, err = f.WriteString(text); err != nil {
		return err
	}

	return nil
}

func setLowest3BitsToZero(n int64) int64 {
	return n & ^int64(7)
}

func TTDLEV(c int64, lev int64) int64 {
	return setLowest3BitsToZero(c) | lev
}

func TTX[A any](c int64, v A) A {
	TTD(c, "", "v1", v)
	return v
}

func TTX2[A any, B any](c int64, v1 A, v2 B) (A, B) {
	TTD(c, "", "v1", v1, "v2", v2)
	return v1, v2
}

func TTX3[A any, B any, C any](c int64, v1 A, v2 B, v3 C) (A, B, C) {
	TTD(c, "", "v1", v1, "v2", v2, "v3", v3)
	return v1, v2, v3
}

func TTD1(rq int64, message string, p1 string, v1 int64) {
	TTD(rq, message, p1, v1)
}

func TTD2(rq int64, message string, p1 string, v1 int64, p2 string, v2 int64) {
	TTD(rq, message, p1, v1, p2, v2)
}

func TTD3(rq int64, message string, p1 string, v1 int64, p2 string, v2 int64) {
	TTD(rq, message, p1, v1, p2, v2)
}

func __TTD(rq int64, message string, pars ...interface{}) {
}

func TTD(rq int64, message string, pars ...interface{}) {
	x := 1
	ptr, fileName, line, isOK := runtime.Caller(x)
	for {
		if filePath(fileName) == "ttd.go" {
			x = x + 1
			ptr, fileName, line, isOK = runtime.Caller(x)
		} else {
			break
		}
	}
	_ = ptr
	_ = isOK
	fl := &fileLine{}
	original := fl
	newfl := fl
	fl.FN = filePath(fileName)
	fl.LN = line
	fl.Path = fileName
	fl.M = message

	m := NewOrderedMap()
	for q := 0; q < len(pars); q = q + 2 {
		m.Set(fmt.Sprintf("%v", pars[q]), fmt.Sprintf("%v", pars[q+1]))
		//m[fmt.Sprintf("%v", pars[x])] = fmt.Sprintf("%v", pars[x+1])
	}
	fl.P = m
	x = x + 1
	for {
		ptr, fileName, line, isOK := runtime.Caller(x)
		_ = ptr
		newfl = &fileLine{}
		if isOK {
			newfl.FN = filePath(fileName)
			newfl.LN = line
			newfl.Path = fileName
		} else {
			break
		}
		fl.Parent = newfl
		fl = newfl
		x = x + 1
	}

	//fmt.Printf("Caller")
	jsonBytes, err := json.Marshal(original)
	_ = err
	//fmt.Print(string(jsonBytes))

	cwd, ok := os.Getwd()
	_ = ok
	fp := filepath.Join(cwd, "..", "working_data", "ttd", strconv.FormatInt(setLowest3BitsToZero(rq), 10)+".txt")
	//fmt.Print(fp)
	appendToFile(fp, string(jsonBytes)+"\n")

	/*fmt.Printf(" %d ", ptr)
	fmt.Printf(fileName)
	fmt.Printf(" %d ", line)
	fmt.Printf(" %d ", isOK)

	for _, par := range pars {
		fmt.Print(par)
	}

	fmt.Printf("\n\n")*/
}
