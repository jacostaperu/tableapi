package tableapi

import (
	"context"
	"embed"
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"gopkg.in/ini.v1"

	adminGUI "github.com/jacostaperu/tableapi.git/web/admin/templs"
)

type Server struct {
	router     *http.ServeMux
	TablesPath string // relative where the program is run
	devMode    bool
	Logger     *Logger
}

//go:embed configSample.conf
var configData embed.FS

//go:embed all:static
var embeddedFiles embed.FS

func NewServer() *Server {
	s := Server{
		router:     http.NewServeMux(),
		TablesPath: "tables",
		devMode:    true,
		Logger:     NewLogger("debug"),
	}

	s.routes()

	return &s
}

func (s *Server) RunDevMode() {
	s.devMode = true
}
func (s *Server) SetTablesPath(tablespath string) {
	s.TablesPath = tablespath
}

func (s *Server) RunProdMode() {
	s.devMode = false
}

func (s *Server) handleFrontEnd() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Debugf("aqui estoi iniciando la pagina web %+v  ", s)

		adminGUI := adminGUI.AdminGUI()

		adminGUI.Render(context.Background(), w)

	}
}

func (s *Server) NewServeMux() *http.ServeMux {
	return s.router
}

func (s *Server) handleStatic() http.HandlerFunc {
	subFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(subFS))

	fileServer = http.StripPrefix("/static", fileServer)
	//fileServer = http.StripPrefix("/static", fileServer)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("aqui estoi iniciando la pagina web static %+v  ", s)

		requestedPath := r.URL.Path
		log.Println("requestedPath embedd: ", requestedPath)
		// Add any extra logic here (e.g., logging)
		log.Printf("Serving embedded file: %s", r.URL.Path)

		// Delegate request to embedded file server
		fileServer.ServeHTTP(w, r)

		//serveEmbeddedFile(w, r)

	}
}

func (s *Server) handleReadTable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tablename := r.PathValue("tablename")
		tabledata, err := s.loadCSV(tablename)
		if err != nil {
			w.Write([]byte("there was an error"))
		}

		log.Println("reading " + tablename)
		log.Println(tabledata)
		readTable := adminGUI.ReadTable(tablename, tabledata)
		readTable.Render(context.Background(), w)

	}
}

func (s *Server) handleEditRow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tablename := r.PathValue("tablename")
		id := r.PathValue("id")
		idx := r.PathValue("idx")
		tabledata, err := s.loadCSV(tablename)
		if err != nil {
			w.Write([]byte("there was an error"))
		}
		// now select the row by ID
		var row []string
		var titles []string
		for _, r := range tabledata {
			if r[0] == id {
				log.Println(r)
				row = append(row, r...)
				titles = append(titles, tabledata[0]...)
			}
		}

		log.Println("reading row for EDIT of " + tablename)
		log.Println("id " + id)
		log.Println("idx " + idx)
		log.Println(row)
		log.Println(titles)

		editRow := adminGUI.EditRow(tablename, idx, titles, row)
		editRow.Render(context.Background(), w)

	}
}

func (s *Server) handleReadRow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tablename := r.PathValue("tablename")
		id := r.PathValue("id")
		idx := r.PathValue("idx")
		tabledata, err := s.loadCSV(tablename)
		if err != nil {
			w.Write([]byte("there was an error"))
		}
		// now select the row by ID
		var row []string
		for _, r := range tabledata {
			if r[0] == id {
				log.Println(r)
				row = append(row, r...)
			}
		}

		log.Println("reading row of " + tablename)
		log.Println("id " + id)
		log.Println("idx " + idx)
		log.Println(row)
		editRow := adminGUI.ReadRow(tablename, idx, row)
		editRow.Render(context.Background(), w)

	}
}

func (s *Server) handleSaveRow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("ieaieaiea", r)
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		tablename := r.PathValue("tablename")
		id := r.PathValue("id")
		idx := r.PathValue("idx")
		tabledata, err := s.loadCSV(tablename)
		log.Println(tabledata)
		if err != nil {
			w.Write([]byte("there was an error"))
		}
		// now select the row by ID
		var row []string
		var titles []string
		for i, crow := range tabledata {
			if crow[0] == id {
				log.Println(r)
				//here I need to update the data
				for j, ftitle := range tabledata[0] {
					if j == 0 {
						continue
					}
					log.Println("updating ", ftitle, " on tabledada[", i, "][", j, "] =", r.FormValue(ftitle), " it was ", crow[j])
					tabledata[i][j] = r.FormValue(ftitle)

				}
				row = append(row, tabledata[i]...)

				break
			}
		}

		//now the tabledata need to persist to disk
		err = s.saveCSV(tablename, tabledata)
		if err != nil {
			log.Println(err)
		}
		log.Println("reading to save row of " + tablename)
		log.Println("id " + id)
		log.Println("idx " + idx)
		log.Println(row)
		log.Println("tabledata", tabledata)
		log.Println(titles)
		editRow := adminGUI.ReadRow(tablename, idx, row)
		editRow.Render(context.Background(), w)

	}
}

// this is to load the table stored in json format
func (s *Server) loadCSV(filename string) ([][]string, error) {
	file, err := os.Open(s.TablesPath + "/" + filename + ".csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}
	return records, nil
}

func (s *Server) saveCSV(filename string, tabledata [][]string) error {
	file, err := os.Create(s.TablesPath + "/" + filename + ".csv")
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Write all records
	log.Println("About to write CSV", tabledata)

	err = writer.WriteAll(tabledata) // Calls Flush internally
	if err != nil {
		return fmt.Errorf("Failed to write data to CSV: %v", err)
	}

	// Optionally flush manually (since WriteAll already flushes)
	writer.Flush()
	return nil
}

func (s *Server) tabledata2GenericTable(tabledata [][]string) GenericTable {
	var records = make([]GenericRecord, 0)
	var titles = make([]string, 0)
	//get titles
	titles = append(titles, tabledata[0]...)

	for i, row := range tabledata {
		if i == 0 {
			continue
		}

		record := GenericRecord{
			ID:     row[0],
			Fields: make(map[string]string),
		}
		for j, r := range row {
			if j == 0 {
				continue
			}
			record.Fields[titles[j]] = r
		}
		records = append(records, record)
	}

	genericTable := GenericTable{
		Records: records,
	}
	return genericTable
}

func (s *Server) LoadConfig() (*ini.File, error) {
	paths := []string{"./tableapi.conf", "/etc/tableapi.conf"}
	var cfg *ini.File
	var err error
	for _, path := range paths {
		if _, err = os.Stat(path); err == nil {
			cfg, err = ini.Load(path)
			if err != nil {
				return nil, fmt.Errorf("failed to load config from %s: %v", path, err)
			}
			fmt.Printf("Loaded config from %s\n", path)
			return cfg, nil
		}
	}
	return nil, fmt.Errorf("config file not found in any expected location:\n   ./tableapi.conf\n   /etc/tableapi.conf")
}

func (s *Server) CreateConfig() error {
	s.Logger.Info("Creating a tableapi.conf")
	data, err := configData.ReadFile("configSample.conf")
	if err != nil {
		return err
	}

	err = os.WriteFile("tableapi.conf", []byte(data), 0644)
	return err
}
