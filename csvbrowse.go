package main

import (
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/zetamatta/go-mbcs"
)

func do_file(fname string, w io.Writer) error {
	r, r_err := os.Open(fname)
	if r_err != nil {
		return r_err
	}
	defer r.Close()

	ansi_all, ansi_all_err := ioutil.ReadAll(r)
	if ansi_all_err != nil {
		return ansi_all_err
	}

	unicode_all, unicode_all_err := mbcs.AtoU(ansi_all)
	if unicode_all_err != nil {
		return unicode_all_err
	}
	csvr := csv.NewReader(strings.NewReader(unicode_all))
	for {
		cols, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			if err == csv.ErrFieldCount {
				goto safe
			}
			if t := err.(*csv.ParseError); t != nil && t.Err == csv.ErrFieldCount {
				goto safe
			}
			return err
		}
	safe:
		fmt.Fprint(w, "<tr>")
		for _, c := range cols {
			fmt.Fprintf(w, "<td nowrap>%s</td>", html.EscapeString(c))
		}
		fmt.Fprintln(w, "</tr>")
	}
}

func main1(files []string, htmlpath string) error {
	w, w_err := os.Create(htmlpath)
	if w_err != nil {
		return w_err
	}
	fmt.Fprintln(w, `<html>`)
	fmt.Fprintln(w, `<head>`)
	fmt.Fprintln(w, `<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />`)
	fmt.Fprintln(w, `</head>`)
	fmt.Fprintln(w, `<body><table border>`)
	defer func() {
		fmt.Fprintln(w, `</table></body></html>`)
		w.Close()
	}()
	for _, fname := range files {
		if err := do_file(fname, w); err != nil {
			return fmt.Errorf("%s: %s", fname, err.Error())
		}
	}
	return nil
}

const htmlpath = "tmp.html"

func main() {
	if err := main1(os.Args[1:], htmlpath); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	cmd1 := exec.Cmd{
		Path: "cmd.exe",
		Args: []string{"/c", "start", htmlpath},
	}
	cmd1.Run()
}