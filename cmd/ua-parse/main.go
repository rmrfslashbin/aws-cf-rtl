package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/rmrfslashbin/aws-cf-rtl/pkg/useragent"
)

const (
	// hard code or pass as arg 1
	datafile = "" // tsv file with user-agent as second column
)

type kv struct {
	Key   string
	Value int
}

func main() {
	var tsvfile string
	if len(os.Args) > 1 {
		tsvfile = os.Args[1]
	} else {
		tsvfile = datafile
	}

	if tsvfile == "" {
		log.Fatal("no data file specified")
	}

	// Sanatize the datafile path.
	fqpn := path.Clean(tsvfile)
	fmt.Println(fqpn)

	family := make(map[string]int)
	osFamily := make(map[string]int)
	devFamily := make(map[string]int)

	// Open the file.
	datafile, err := os.Open(fqpn)
	if err != nil {
		log.Fatal(err)
	}
	defer datafile.Close()

	tsv := csv.NewReader(datafile)
	tsv.Comma = '\t'

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "family\tmajor\tminor\tpatch\tos_family\tos_major\tos_minor\tos_patch\tos_patch_minor\tdevice_family\tdevice_brand\tdevice_model")

	for {
		record, err := tsv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// ip := strings.TrimSpace(record[0])
		ua := strings.TrimSpace(record[1])

		uaRecord := useragent.Parse(ua)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			uaRecord.UAFamily,
			uaRecord.UAMajor,
			uaRecord.UAMinor,
			uaRecord.UAPatch,
			uaRecord.UAOSFamily,
			uaRecord.UAOSMajor,
			uaRecord.UAOSMinor,
			uaRecord.UAOSPatch,
			uaRecord.UAOSPatchMinor,
			uaRecord.UADeviceFamily,
			uaRecord.UADeviceBrand,
			uaRecord.UADeviceModel)

		family[uaRecord.UAFamily]++
		osFamily[uaRecord.UAOSFamily]++
		devFamily[uaRecord.UADeviceFamily]++
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "family\tcount")

	for _, kv := range srt(&family) {
		fmt.Fprintf(w, "%s\t%d\n", kv.Key, kv.Value)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "os\tcount")

	for _, kv := range srt(&osFamily) {
		fmt.Fprintf(w, "%s\t%d\n", kv.Key, kv.Value)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "device\tcount")

	for _, kv := range srt(&devFamily) {
		fmt.Fprintf(w, "%s\t%d\n", kv.Key, kv.Value)
	}

	w.Flush()

}

func srt(m *map[string]int) []kv {
	var ss []kv
	for k, v := range *m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return ss
}
