// The rjson command reads rjson or JSON from its standard input and
// writes rjson (or JSON if the "-j" flag is given) to its standard output.
// The "-indent" flag specifies an indent string; if this is empty,
// a compact output will be used; otherwise the result will be
// pretty-printed.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"github.com/egonelbre/spexs2/vendor/rjson"
	"os"
)

var jsonOut = flag.Bool("j", false, "produce JSON output")
var indent = flag.String("indent", "\t", "indentation string; empty implies compact output")

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "rjson reads from stdin only")
		os.Exit(2)
	}
	var marshal func(interface{}) ([]byte, error)
	if *jsonOut {
		if *indent == "" {
			marshal = json.Marshal
		} else {
			marshal = func(v interface{}) ([]byte, error) {
				return json.MarshalIndent(v, "", *indent)
			}
		}
	} else {
		if *indent == "" {
			marshal = func(v interface{}) ([]byte, error) {
				data, err := rjson.MarshalIndent(v, "", "")
				if err != nil {
					return data, err
				}
				// TODO: use Compact only when it compacts
				// identifiers too.
				var b bytes.Buffer
				if err := rjson.Compact(&b, data); err != nil {
					return nil, err
				}
				return b.Bytes(), nil
			}
		} else {
			marshal = func(v interface{}) ([]byte, error) {
				return rjson.MarshalIndent(v, "", *indent)
			}
		}
	}
	d := rjson.NewDecoder(os.Stdin)
	d.UseNumber()
	w := bufio.NewWriter(os.Stdout)

	err := stream(w, d, marshal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rjson: %v", err)
		os.Exit(1)
	}
}

func stream(w *bufio.Writer, d *rjson.Decoder, marshal func(v interface{}) ([]byte, error)) error {
	var v interface{}
	defer w.Flush()
	for {
		err := d.Decode(&v)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decode: %v", err)
		}
		data, err := marshal(v)
		if err != nil {
			return fmt.Errorf("marshal: %v", err)
		}
		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("write: %v", err)
		}
		v = nil
		w.WriteByte('\n')
	}
	return nil
}
