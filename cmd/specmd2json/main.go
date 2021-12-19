/*
Copyright © 2021 Z5Labs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/z5labs/bbcos/cmd/specmd2json/parser"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var (
	specFileUrl = flag.String("spec", "", "URL of markdown spec file")
	output      = flag.String("output", "", "Filename to output to.")
)

func main() {
	flag.Parse()
	*specFileUrl = strings.TrimSpace(*specFileUrl)
	if *specFileUrl == "" {
		log.Fatal("specmd2json: spec file url must be non empty")
	}

	ext := path.Ext(*specFileUrl)
	if ext != ".md" {
		log.Fatal("specmd2json: spec file must be a markdown file")
	}

	resp, err := http.Get(*specFileUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	spec := parser.Parse(resp.Body)

	var w io.Writer = os.Stdout
	out := strings.TrimSpace(*output)
	if len(out) > 0 {
		w, err = os.Create(out)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = writeTo(spec, w)
	if err != nil {
		log.Fatal(err)
	}
}

func writeTo(v *viper.Viper, w io.Writer) error {
	fs := afero.NewMemMapFs()
	v.SetFs(fs)
	v.SetConfigFile("test.json")
	err := v.WriteConfig()
	if err != nil {
		return err
	}

	f, err := fs.Open("test.json")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	return err
}
