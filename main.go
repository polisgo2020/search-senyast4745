package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/urfave/cli/v2"

	"github.com/polisgo2020/search-senyast4745/index"
	"github.com/polisgo2020/search-senyast4745/util"
)

func main() {

	var err error

	log.Println(`
	 ___ _   ___     ______  _____    _    ____   ____ _   _ 
	|_ _| \ | \ \   / / ___|| ____|  / \  |  _ \ / ___| | | |
	 | ||  \| |\ \ / /\___ \|  _|   / _ \ | |_) | |   | |_| |
	 | || |\  | \ V /  ___) | |___ / ___ \|  _ <| |___|  _  |
	|___|_| \_|  \_/  |____/|_____/_/   \_\_| \_\\____|_| |_|

	`)

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
	portFlag := &cli.StringFlag{
		Aliases:     []string{"p"},
		Name:        "port",
		Usage:       "Network interface",
		DefaultText: "8888",
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
			},
			Action: build,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search over the index",
			Flags: []cli.Flag{
				indexFileFlag,
				portFlag,
				debugFlag,
			},
			Action: search,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Printf("Fatal with %q error while starting command line app", err)
	}
}

func build(c *cli.Context) error {

	log.Println("msg", "build run", "index file", c.String("index"),
		"source folder", c.String("sources"))

	if err := checkFlags(c, "index", "sources"); err != nil {
		log.Println("error", err, "context flags", c.FlagNames(),
			"msg", "error while checking context")
		return nil
	}
	if allFiles, err := filePathWalkDir(c.String("sources")); err != nil {
		log.Println("error", err,
			"msg", fmt.Sprintf("can not read files list from directory: %s", c.String("sources")))
	} else {
		log.Println("msg", "folder parsed", "files", fmt.Sprintf("%+v", allFiles))

		m := collectWordData(allFiles)

		log.Println("msg", "index built")
		if err := collectAndWriteMap(m, c.String("index")); err != nil {
			log.Println("error", err,
				"msg", fmt.Sprintf("can not save data to file with name: %s", c.String("index")))
			return nil
		} else {
			log.Println("msg", "index saved")
		}
	}

	log.Println("msg", "done")

	return nil
}

func collectWordData(fileNames []string) *index.Index {
	m := index.NewIndex()

	m.OpenApplyAndListenChannel(func(wg *sync.WaitGroup) {
		for i := range fileNames {
			wg.Add(1)
			go readFileByWords(wg, m, fileNames[i])
		}
		log.Println("msg", fmt.Sprintf("goroutine count %d", len(fileNames)))
	})

	return m
}

func collectAndWriteMap(ind *index.Index, indexFile string) error {
	log.Println("msg", "writing index to file in csv format", "file", indexFile, "index", ind,
		"index length", len(ind.Data))
	recordFile, _ := os.Create(indexFile)
	return ind.ToFile(index.NewCsvEncoder(recordFile))
}

type FileResponse struct {
	Filename string
	Count    int
	Spacing  int
}

func search(c *cli.Context) error {

	log.Println("msg", "search run", "index file", c.String("index"),
		"server port", c.String("port"))

	if err := checkFlags(c, "index"); err != nil {
		log.Println("error", err, "context flags", c.FlagNames(),
			"msg", "error while checking context")
		return nil
	}

	wapp, err := NewApp(c.String("port"))

	if err != nil {
		log.Println("error", err, "msg", "error while creating web application")
		return nil
	}

	data, err := readCSVFile(c.String("index"))
	if err != nil {
		log.Println("error", err,
			"msg", fmt.Sprintf("Couldn't open or read the csv file %s", c.String("path")))
		return nil
	}

	r := wapp.Mux

	r.Post("/", func(w http.ResponseWriter, req *http.Request) {
		searchWords := req.FormValue("search")
		log.Println("msg", searchWords)
		var inputWords []string
		for _, word := range strings.Split(searchWords, ",") {
			util.CleanUserInput(word, func(input string) {
				inputWords = append(inputWords, input)
			})
		}
		log.Println("msg", inputWords)
		if len(inputWords) == 0 {
			log.Println("error", nil,
				"msg", "Incorrect search words", "input", c.String("search-words"))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var resp []FileResponse
		for k, v := range data.Search(inputWords) {
			resp = append(resp, FileResponse{
				Filename: k,
				Count:    v.Path,
				Spacing:  v.Weight,
			})
		}
		log.Println("msg", fmt.Sprintf("resp %+v", resp))

		rawData, err := json.Marshal(resp)
		if err != nil {
			log.Printf("error %s while marshalling data %+v to json\n", err, resp)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		if _, err = fmt.Fprint(w, string(rawData)); err != nil {
			log.Printf("error %s while writing data %s do json\n", err, string(rawData))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":"+wapp.Port, r); err != nil {
		log.Println("error", err)
	}

	return nil
}

func readCSVFile(filePath string) (*index.Index, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	data := index.NewIndex()
	decoder := index.NewCsvDecoder(csvFile)
	err = data.FromFile(decoder)
	return data, err
}

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

func readFileByWords(wg *sync.WaitGroup, ind *index.Index, fn string) {
	defer wg.Done()
	log.Println("msg", "goroutine start", "filename", fn, "goroutine id", goid())
	file, err := os.Open(fn)
	if err != nil {
		log.Println("error", err,
			"msg", fmt.Sprintf("can not open file %s", fn), "filename", fn)
		return
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()

	ind.MapAndCleanWords(file, fn)
	log.Println("msg", "goroutine normal end", "goroutine id", goid())
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

type App struct {
	Mux  *chi.Mux
	Port string
}

func NewApp(port string) (*App, error) {
	r := chi.NewMux()
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Timeout(100 * time.Millisecond))
	return &App{Mux: r, Port: port}, nil
}
