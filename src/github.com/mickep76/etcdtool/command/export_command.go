package command

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"
	"strconv"
)

// NewExportCommand returns data from export.
func NewExportCommand() cli.Command {
	return cli.Command{
		Name:  "export",
		Usage: "export a directory",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "sort, s", Usage: "returns result in sorted order"},
			cli.StringFlag{Name: "format, f", EnvVar: "ETCDTOOL_FORMAT", Value: "JSON", Usage: "Data serialization format YAML, TOML or JSON"},
			cli.StringFlag{Name: "output, o", Value: "", Usage: "Output file"},
			cli.BoolFlag{Name: "num-infer-list", Usage: "returns result without extra levels of arrays"},
		},
		Action: func(c *cli.Context) error {
			exportCommandFunc(c)
			return nil
		},
	}
}

// exportCommandFunc exports data as either JSON, YAML or TOML.
func exportCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		fatal("You need to specify directory")
	}
	dir := c.Args()[0]

	// Remove trailing slash.
	if dir != "/" {
		dir = strings.TrimRight(dir, "/")
	}
	infof("Using dir: %s", dir)

	// Load configuration file.
	e := loadConfig(c)

	// New dir API.
	ki := newKeyAPI(e)

	sort := c.Bool("sort")

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		fatal(err.Error())
	}

	exportFunc(dir, sort, c.String("output"), f, c, ki)
}

// exportCommandFunc exports data as either JSON, YAML or TOML.
func exportFunc(dir string, sort bool, file string, f iodatafmt.DataFmt, c *cli.Context, ki client.KeysAPI) {
	ctx, cancel := contextWithCommandTimeout(c)
	resp, err := ki.Get(ctx, dir, &client.GetOptions{Sort: sort, Recursive: true})
	cancel()
	if err != nil {
		fatal(err.Error())
	}

	m := etcdmap.Map(resp.Node)
	if c.Bool("num-infer-list") {
		m1 := removeExtraNumbersLevels(m)
		value, ok := m1.(map[string]interface{})
		if ok {
			m = value
		}
	}

	// Export and write output.
	if file != "" {
		iodatafmt.Write(file, m, f)
	} else {
		iodatafmt.Print(m, f)
	}
}
/*
	Remove extra levels of numbers created in etcd.
 */
func removeExtraNumbersLevels(etcdmapObject interface{}) interface{} {

	var result map[string]interface{} = make(map[string]interface{})
	// START TRAVERSE MAP
	switch etcdmapObject.(type) {
	case map[string]interface{}: // map {string, K} case

		for k, v := range etcdmapObject.(map[string]interface{}) {
			// START TRAVERSE VALUES TYPE
			switch v.(type) {

			case map[string]interface{}:

				allKeyNumbers := checkAllKeysAreNumbers(v)

				if allKeyNumbers {
					// traverse the values to create an array
					// and removeExtraNumbersLevels in the subsequent levels
					var results []interface{}
					for _, allKeyNumbersValue := range v.(map[string]interface{}) {
						// create temporal map with a fake top level to be in accordance
						// with the logic of the function
						dummy := make(map[string]interface{})
						fake_key := "flatten_fake_key"
						dummy[fake_key] = allKeyNumbersValue
						// flat the temporal map
						flatten := removeExtraNumbersLevels(dummy)
						// set the results depends on the type of the map returned
						value, ok := flatten.(map[string]interface{})
						if ok {
							results = append(results, value[fake_key])
						} else {
							results = append(results, flatten)
						}
					}
					// set the processed subkeys to the result map
					if len(results) == 0 {
						c := []string{}
						result[k] = c
					} else {
						result[k] = results
					}
				} else {
					// set a normal key and removeExtraNumbersLevels in the subsequent levels
					result[k] = removeExtraNumbersLevels(v)
				}
			default:
				// process a normal value
				result[k] = v
			}
		}
	case string:
		// return a normal value
		return etcdmapObject
	}

	return result
}

func checkAllKeysAreNumbers(numbersMap interface{}) bool {

	allKeyNumbers := true
	for k, _ := range numbersMap.(map[string]interface{}) {
		_, err := strconv.Atoi(k)
		if err != nil {
			allKeyNumbers = false
			break
		}
	}
	return allKeyNumbers

}

