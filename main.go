package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/polisgo2020/search-senyast4745/index"
	"github.com/polisgo2020/search-senyast4745/log"
	"github.com/polisgo2020/search-senyast4745/util"
	"github.com/urfave/cli/v2"
)

var wapp *App

func main() {

	//os.Mkdir("logs")

	var err error
	wapp, err = NewApp()

	if err != nil {
		fmt.Printf("Error %e while starting application", err)
		return
	}

	fmt.Print("    ___ _   ___     ______  _____    _    ____   ____ _   _\n" +
		"   |_ _| \\ | \\ \\   / / ___|| ____|  / \\  |  _ \\ / ___| | | |\n" +
		"    | ||  \\| |\\ \\ / /\\___ \\|  _|   / _ \\ | |_) | |   | |_| |\n" +
		"    | || |\\  | \\ V /  ___) | |___ / ___ \\|  _ <| |___|  _  |\n" +
		"   |___|_| \\_|  \\_/  |____/|_____/_/   \\_\\_| \\_\\\\____|_| |_|\n\n")

	app := cli.NewApp()

	app.Version = "0.0.1"
	app.Authors = []*cli.Author{{Name: "Arseny Druzhinin", Email: "senyasdt4745@gmail.com"}}
	app.Name = "Search index"
	app.Usage = "generate index from text files and search over them"

	indexFileFlag := &cli.StringFlag{
		Aliases:     []string{"i"},
		Name:        "index",
		Usage:       "Index file",
		DefaultText: "output/final.csv",
	}

	sourcesFlag := &cli.StringFlag{
		Aliases:  []string{"s"},
		Name:     "sources, s",
		Usage:    "Files to index",
		Required: true,
	}

	searchFlag := &cli.StringFlag{
		Aliases:  []string{"sw"},
		Name:     "search-word, sw",
		Usage:    "Search words separated by comma",
		Required: true,
	}

	logFolderFlag := &cli.BoolFlag{
		Name:  "log",
		Usage: "Turn on logging to files",
	}

	debugFlag := &cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"d"},
		Usage:   "Turn on debug mode",
	}

	app.Commands = []*cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Build search index",
			Flags: []cli.Flag{
				indexFileFlag,
				sourcesFlag,
				debugFlag,
				logFolderFlag,
			},
			Action: build,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search over the index",
			Flags: []cli.Flag{
				indexFileFlag,
				searchFlag,
				debugFlag,
				logFolderFlag,
			},
			Action: search,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("Fatal with %e error while starting command line app", err)
	}
}

func build(c *cli.Context) error {

	log.Debug("msg", "build run", "index file", c.String("index"),
		"source folder", c.String("sources"))

	if err := checkFlags(c, "index", "sources"); err != nil {
		log.Error("error", err, "context flags", c.FlagNames(),
			"msg", "error while checking context")
		return nil
	}
	if allFiles, err := filePathWalkDir(c.String("sources")); err != nil {
		log.Error("error", err,
			"msg", fmt.Sprintf("can not read files list from directory: %s", c.String("sources")))
	} else {
		log.Debug("msg", "folder parsed", "files", fmt.Sprintf("%+v", allFiles))

		m := collectWordData(allFiles)

		log.Debug("msg", "index built")
		if err := collectAndWriteMap(m, c.String("index")); err != nil {
			log.Error("error", err,
				"msg", fmt.Sprintf("can not save data to file with name: %s", c.String("index")))
			return nil
		} else {
			log.Debug("msg", "index saved")
		}
	}

	log.Info("msg", "done")

	return nil
}

func collectWordData(fileNames []string) *index.Index {
	m := index.NewIndex()

	m.OpenApplyAndListenChannel(func(wg *sync.WaitGroup) {

		for i := range fileNames {
			wg.Add(1)
			go readFileByWords(wg, m, fileNames[i])
		}
		log.Debug("msg", fmt.Sprintf("goroutine count %d", len(fileNames)))
	})

	return m
}

func collectAndWriteMap(ind *index.Index, indexFile string) error {
	log.Debug("msg", "writing index to file in csv format", "file", indexFile, "index", ind,
		"index length", len(ind.Data))
	recordFile, _ := os.Create(indexFile)
	w := csv.NewWriter(recordFile)
	var count int
	for k, v := range ind.Data {
		t, err := json.Marshal(v)
		if err != nil {
			log.Error("error", err,
				"msg", fmt.Sprintf("can not create json from obj %+v \n", &v), "obj", &v)
		}
		err = w.Write([]string{k, string(t)})
		if err != nil {
			log.Error("error", err,
				"msg", fmt.Sprintf("can not save record %s,%s \n", k, t), "file", indexFile, "word", k,
				"filestr", string(t))
		}
		count++
		if count > 10 {
			w.Flush()
			log.Debug("msg", "flush writer", "writer", w)
			count = 0
		}
	}
	return nil
}

type FileResponse struct {
	Filename string
	Count    int
	Spacing  int
}

func search(c *cli.Context) error {

	log.Debug("msg", "search run", "index file", c.String("index"),
		"search words", c.String("search-word"), "server port", wapp.Port)
	if err := checkFlags(c, "index"); err != nil {
		log.Error("error", err, "context flags", c.FlagNames(),
			"msg", "error while checking context")
		return nil
	}

	data, err := readCSVFile(c.String("index"))
	if err != nil {
		log.Error("error", err,
			"msg", fmt.Sprintf("Couldn't open or read the csv file %s", c.String("path")))
		return nil
	}

	r := wapp.Mux

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		searchWords := r.FormValue("search")
		log.Info("msg", searchWords)
		inputWords := make([]string, 0)
		for _, word := range strings.Split(searchWords, ",") {
			util.CleanUserInput(word, func(input string) {
				inputWords = append(inputWords, input)
			})
		}
		log.Info("msg", inputWords)
		if len(inputWords) == 0 {
			log.Error("error", nil,
				"msg", "Incorrect search words", "input", c.String("search-words"))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var resp []FileResponse
		for k, v := range getCorrectFiles(data, inputWords) {
			resp = append(resp, FileResponse{
				Filename: k,
				Count:    v.Path,
				Spacing:  v.Weight,
			})
		}
		log.Info("msg", fmt.Sprintf("resp %+v", resp))

		rawData, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		if _, err = fmt.Fprint(w, string(rawData)); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":"+wapp.Port, r); err != nil {
		log.Error("error", err)
	}

	return nil
}

func getCorrectFiles(m *index.Index, searchWords []string) map[string]*index.Data {
	data := index.NewIndex()
	for i := range searchWords {
		tmp := m.Data[searchWords[i]]
		if len(tmp) != 0 {
			data.Data[searchWords[i]] = tmp
		}
	}
	log.Debug("msg", "start search in index")
	return data.Search(searchWords)
}

func readCSVFile(filePath string) (*index.Index, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(csvFile)
	data := index.NewIndex()
	var errCount int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("error", err,
				"msg", "can not read csv line")
			errCount++
			if errCount > 100 {
				return nil, err
			}
			continue
		}
		log.Debug("msg", "reading data from csv", "data", record)
		var tmp []*index.FileStruct
		if json.Unmarshal([]byte(record[1]), &tmp) != nil {
			log.Error("error", err,
				"msg", fmt.Sprintf("can not parse json data %s \n", record[1]), "data", record[1])
			continue
		}
		data.Data[record[0]] = tmp
	}
	return data, nil
}

// FilePathWalkDir bypasses the given director and returns a list of all files in this folder
// and returns an error if it is not possible to access the folder
func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return errors.New("")
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ReadFileByWords reads the given file by words and returns an array of the layer
//or an error if it is impossible to open or read the file
func readFileByWords(wg *sync.WaitGroup, ind *index.Index, fn string) {
	defer wg.Done()
	log.Debug("msg", "goroutine start", "filename", fn, "goroutine id", goid())
	file, err := os.Open(fn)
	if err != nil {
		log.Error("error", err,
			"msg", fmt.Sprintf("can not open file %s", fn), "filename", fn)
		return
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()

	ind.MapAndCleanWords(file, fn)
	log.Debug("msg", "goroutine normal end", "goroutine id", goid())
	return
}

func goid() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	return idField
}

func checkFlags(c *cli.Context, str ...string) error {
	for _, flag := range str {
		if c.String(flag) == "" {
			return errors.New(fmt.Sprintf("empty flag %s", flag))
		}
	}
	return nil
}

type appYamlConfig struct {
	Logger []log.Config `yaml:"logger"`
	Server struct {
		Port    string `yaml:"port"`
		Timeout int    `yaml:"timeout"`
	}
}

type App struct {
	Mux  *chi.Mux
	Port string
}

func NewAppFromConfig(conf *appYamlConfig) (*App, error) {

	r := chi.NewRouter()
	r.Use(func(handler http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				log.Debug("msg", "response", "status", ww.Status(), "written", ww.BytesWritten(), "time", time.Since(t1))
			}()

			handler.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	})
	r.Use(middleware.Timeout(time.Duration(conf.Server.Timeout) * time.Millisecond))

	log.GetLogger(conf.Logger...)
	return &App{Mux: r, Port: conf.Server.Port}, nil
}

func NewDefaultApp() (*App, error) {
	r := chi.NewMux()
	log.GetLogger()
	return &App{Mux: r, Port: "8080"}, nil
}

func NewAppFromConfigFile(configFile string) (*App, error) {
	if _, err := os.Stat(configFile); err == nil {
		dat, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}
		conf := appYamlConfig{}
		if err = yaml.Unmarshal(dat, &conf); err != nil {
			return nil, err
		}

		return NewAppFromConfig(&conf)
	} else {
		return NewDefaultApp()
	}
}

func NewApp() (*App, error) {
	return NewAppFromConfigFile("config.yml")
}
