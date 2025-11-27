package tableapi

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

// in‑memory store (replace with DB later)
var (
	store = GenericTable{}
	mu    sync.RWMutex
)

// GET /records  -> return all
func (s *Server) handleRecords() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:

			mu.RLock()
			//here load the table
			tablename := r.PathValue("tablename")

			tabledata, err := s.loadCSV(tablename)

			titles := tabledata[0]

			//here the table will be filtered if any of the columns in the table comes as parameter
			filter := make(map[string]string)

			// Parse URL query parameters
			queryParams := r.URL.Query()

			for _, title := range titles {
				value := queryParams.Get(title)
				if value != "" {
					filter[title] = value
					log.Printf("Parameter '%s' value '%s'\n", title, value)
				} else {
					// Parameter "XXX" found, print its value
					log.Printf("Parameter '%s' not found or empty", title)
				}
			}
			result := [][]string{tabledata[0]}

			if len(filter) > 0 {
				//
				//
				for _, row := range tabledata {
					for i, title := range titles {
						if filtervalue, exists := filter[title]; exists {
							if row[i] == filtervalue {
								result = append(result, row)
								continue
							}

						}
					}

				}
			} else {
				result = tabledata
			}

			//preareResponse
			genericTable := s.tabledata2GenericTable(result)
			if err != nil {
				http.Error(w, "error loading table: "+err.Error(), http.StatusMethodNotAllowed)
				return
			}

			defer mu.RUnlock()
			json.NewEncoder(w).Encode(genericTable)

		case http.MethodPost:
			var incoming GenericTable
			if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}

			mu.Lock()
			store.Records = append(store.Records, incoming.Records...)
			mu.Unlock()

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(incoming)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// PATCH /records/{id} -> partial update of Fields
// here an example
/*
 * curl -X PATCH "http://localhost:8087/records/PIN_Table/sthascdloaaee" \
   -H "Content-Type: application/json" \
   -d '{
     "records": [
       {
         "id": "sthascdloaaee",
         "fields": {
           "CLI": "441234567890 ahtsha htshtsh",
           "Customer_Name": "Demo Org B",
           "Customer_Support_Class": "Standard",
           "PIN": "87465",
           "Products": "SD-wan\nWebex Calling"
         }
       }
     ]
   }'
*/

func (s *Server) handleRecordID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Println("PATCH operation ")
		if r.Method != http.MethodPatch {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tablename := r.PathValue("tablename")
		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}
		log.Println(tablename, id)
		// expected body: {"records":[{"id":"sameID","fields":{...}}]}
		var incoming GenericTable
		if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if len(incoming.Records) == 0 {
			http.Error(w, "no records in payload", http.StatusBadRequest)
			return
		}
		patch := incoming.Records[0]

		mu.Lock()
		defer mu.Unlock()

		tabledata, err := s.loadCSV(tablename)
		if err != nil {
			http.Error(w, "error reading table", http.StatusBadRequest)
		}

		titles := buildTitles(tabledata[0])

		for i, record := range tabledata {
			if record[0] == id {
				//found
				log.Print("found", i)

				for k, v := range patch.Fields {
					log.Printf("Updating %s to %s \n", k, v)
					tabledata[i][titles[k]] = strings.ReplaceAll(v, "\n", `\n`)
				}
			}
		}
		err = s.saveCSV(tablename, tabledata)
		if err != nil {
			http.Error(w, "error writing table", http.StatusBadRequest)
		}

	}

}

func getFieldNames(i interface{}) []string {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	var fields []string
	for i := 0; i < val.NumField(); i++ {
		// only export fields have names accessible
		field := typ.Field(i)
		if field.PkgPath == "" { // check if exported
			fields = append(fields, field.Name)
		}
	}
	return fields
}

func buildTitles(titlesarr []string) map[string]int {
	titles := make(map[string]int)
	for i, name := range titlesarr {
		titles[name] = i
	}
	return titles

}
