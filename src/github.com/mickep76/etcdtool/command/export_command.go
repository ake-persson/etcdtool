package command

import (
	"strings"
	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"
	"strconv"
	"sort"
	"regexp"
)

const flatten_fake_key = "flatten_fake_key"
const num_infer_list_flag = "num-infer-list"
const infer_types_flag = "infer-types"
const keep_format_path_flag = "keep-format-path"

// NewExportCommand returns data from export.
func NewExportCommand() cli.Command {
	return cli.Command{
		Name:  "export",
		Usage: "export a directory",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "sort, s", Usage: "returns result in sorted order"},
			cli.StringFlag{Name: "format, f", EnvVar: "ETCDTOOL_FORMAT", Value: "JSON", Usage: "Data serialization format YAML, TOML or JSON"},
			cli.StringFlag{Name: "output, o", Value: "", Usage: "Output file"},
			cli.BoolFlag{Name: num_infer_list_flag, Usage: "returns result without extra levels of arrays"},
			cli.BoolFlag{Name: infer_types_flag, Usage: "convert to original type if conversion is possible"},
			cli.StringSliceFlag{ Name: keep_format_path_flag, Usage: "set one or more paths (allow regex) to keep the string format. Each field or level should be separated by '.'"},
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

var keep_formatted_paths []*regexp.Regexp


// exportCommandFunc exports data as either JSON, YAML or TOML.
func exportFunc(dir string, sort bool, file string, f iodatafmt.DataFmt, c *cli.Context, ki client.KeysAPI) {
	ctx, cancel := contextWithCommandTimeout(c)
	resp, err := ki.Get(ctx, dir, &client.GetOptions{Sort: sort, Recursive: true})
	cancel()
	if err != nil {
		fatal(err.Error())
	}

	m := etcdmap.Map(resp.Node)
	if c.Bool(num_infer_list_flag) || c.Bool(infer_types_flag) {
		if c.StringSlice(keep_format_path_flag) != nil {
			for _,path := range(c.StringSlice(keep_format_path_flag)){
				keep_formatted_paths =append(keep_formatted_paths,regexp.MustCompile(path))
			}
		}

		m1 := removeExtraNumbersLevels(m, c.Bool(num_infer_list_flag), c.Bool(infer_types_flag), "")
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

// Remove extra levels of numbers created in etcd and infer numbers
func removeExtraNumbersLevels(etcdmapObject interface{}, numInferList bool, inferTypes bool, path string) interface{} {

	var result map[string]interface{} = make(map[string]interface{})
	// TRAVERSE MAP
	switch etcdmapObject.(type) {
	case map[string]interface{}: // map {string, K} case

		for k, v := range etcdmapObject.(map[string]interface{}) {
			var path_aux string = k
			if len(path) > 0  {
				if k == flatten_fake_key {
					path_aux = path
				}else {
					path_aux = path + "." + k
				}
			}

			// TRAVERSE VALUES TYPE
			switch v.(type) {
			case map[string]interface{}:
				if numInferList && checkAllKeysAreNumbers(v) {
					// traverse the values to create an array
					// and removeExtraNumbersLevels in the subsequent levels
					var results []interface{}

					value, ok := v.(map[string]interface{})
					if ok {
						results = extractArrayFromFirstLevel(value,numInferList,numInferList, path_aux)
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
					result[k] = removeExtraNumbersLevels(v, numInferList, inferTypes, path_aux)
				}
			default:
				assignValue(result, k, v, inferTypes,path_aux)
			}
		}
	case string:
		// return a normal value
		return etcdmapObject
	}

	return result
}

func extractArrayFromFirstLevel(originalMap map[string]interface{}, numInferList bool, inferTypes bool, path string) []interface{} {
	// sort the keys to ensure the list will be in order
	keys := make([]int, 0)
	for k, _ := range originalMap {
		parsedInt , err := strconv.Atoi(k)
		if err != nil{
			fatal(err.Error())
		}else{
			keys = append(keys, parsedInt)
		}
	}
	sort.Ints(keys)

	// process the map and extract the first level to build the array
	var results []interface{}
	for  _, k := range keys {
		allKeyNumbersValue := originalMap[strconv.Itoa(k)]
		// create temporal map with a fake top level to be in accordance
		// with the logic of the function
		temporal_map := make(map[string]interface{})
		temporal_map[flatten_fake_key] = allKeyNumbersValue
		// flat the temporal map
		flatten := removeExtraNumbersLevels(temporal_map, numInferList, inferTypes, path)
		// set the results depends on the type of the map returned
		value, ok := flatten.(map[string]interface{})
		if ok {
			results = append(results, value[flatten_fake_key])
		} else {
			results = append(results, flatten)
		}
	}
	return results
}


func assignValue(result map[string]interface{}, key string, value interface{}, inferTypes bool, path string) {
	isString, ok := value.(string)
	if ok && inferTypes {

		var keep_original_format bool=false
		for _,reg := range(keep_formatted_paths){
			if reg.MatchString(path){
				keep_original_format = true
				break
			}
		}

		if keep_original_format {
			result[key]=isString
		}else {

			// process a normal value
			val, err := strconv.Atoi(isString)
			if err == nil {
				result[key] = val
			} else {
				val, err := strconv.ParseFloat(isString, 64)
				if err == nil {
					result[key] = val
				} else {
					val, err := strconv.ParseBool(isString)
					if err == nil {
						result[key] = val
					} else {
						result[key] = isString

					}
				}
			}
		}
	} else {
		result[key] = value
	}

}

func checkAllKeysAreNumbers(numbersMap interface{}) bool {

    // zero length map shouldn't be regarded as an array
    if (len(numbersMap.(map[string]interface{})) == 0) {
        return false
    }

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

