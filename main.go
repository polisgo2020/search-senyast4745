package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/polisgo2020/search-senyast4745/config"
	"github.com/polisgo2020/search-senyast4745/web"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/urfave/cli/v2"

	"github.com/polisgo2020/search-senyast4745/index"
	"github.com/polisgo2020/search-senyast4745/util"
)

func main() {

	var err error

	initLogger(config.Load())

	log.Info().Msg(`
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

	app.Commands = []*cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Build search index",
			Flags: []cli.Flag{
				indexFileFlag,
				sourcesFlag,
			},
			Action: build,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search over the index",
			Flags: []cli.Flag{
				indexFileFlag,
			},
			Action: search,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Err(err).Msg("Fatal while starting command line app")
	}
}

func initLogger(c *config.Config) {
	logLvl, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		logLvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLvl)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		log.Debug().
			Str("method", r.Method).
			Str("remote", r.RemoteAddr).
			Str("path", r.URL.Path).
			Int("duration", int(time.Since(start))).
			Msgf("Called url %s", r.URL.Path)
	})

}

func build(c *cli.Context) error {

	log.Info().Msg("build mode run")

	log.Debug().
		Str("index file", c.String("index")).
		Str("source folder", c.String("sources")).
		Msg("build run")

	if err := checkFlags(c, "index", "sources"); err != nil {
		log.Err(err).Strs("context flags", c.FlagNames()).Msg("error while checking context")
		return nil
	}
	if allFiles, err := filePathWalkDir(c.String("sources")); err != nil {
		log.Err(err).Str(" directory", c.String("sources")).Msg("can not read files list")
	} else {
		log.Debug().Strs("files", allFiles).Msg("folder parsed")

		m := collectWordData(allFiles)

		log.Debug().Interface("index", m).Msg("index built")
		if err := collectAndWriteMap(m, c.String("index")); err != nil {
			log.Err(err).Str("filename", c.String("index")).Msg("can not save data to file")
			return nil
		} else {
			log.Debug().Msg("index saved")
		}
	}

	log.Debug().Msg("build done")

	return nil
}

func collectWordData(fileNames []string) *index.Index {
	m := index.NewIndex()

	m.OpenApplyAndListenChannel(func(wg *sync.WaitGroup) {
		for i := range fileNames {
			wg.Add(1)
			go readFileByWords(wg, m, fileNames[i])
		}
		log.Debug().Msg(fmt.Sprintf("goroutine count %d", len(fileNames)))
	})

	return m
}

func collectAndWriteMap(ind *index.Index, indexFile string) error {
	log.Info().Str("file", indexFile).Int("index length", len(ind.Data)).
		Msg("writing index to file in csv format")
	recordFile, _ := os.Create(indexFile)
	return ind.ToFile(index.NewCsvEncoder(recordFile))
}

type FileResponse struct {
	Filename string
	Count    int
	Spacing  int
}

func search(c *cli.Context) error {

	log.Info().Msg("search mode run")

	log.Debug().Str("index file", c.String("index")).Interface("config", config.Load()).
		Msg("search run")

	if err := checkFlags(c, "index"); err != nil {
		log.Err(err).Strs("context flags", c.FlagNames()).Msg("error while checking context")
		return nil
	}

	wapp, err := web.NewApp(config.Load(), logMiddleware)

	if err != nil {
		log.Err(err).Msg("error while creating web application")
		return nil
	}

	data, err := readCSVFile(c.String("index"))
	if err != nil {
		log.Err(err).Str("file", c.String("path")).Msg("Couldn't open or read the csv file ")
		return nil
	}

	r := wapp.Mux

	r.Post("/", func(w http.ResponseWriter, req *http.Request) {
		searchWords := req.FormValue("search")
		log.Info().Str("search phrase", searchWords).Msg("start search")
		var inputWords []string
		for _, word := range strings.Split(searchWords, " ") {
			util.CleanUserInput(word, func(input string) {
				inputWords = append(inputWords, input)
			})
		}
		log.Debug().Msgf("clean input: %+v", inputWords)
		if len(inputWords) == 0 {
			log.Err(nil).Str("input", c.String("search-words")).Msg("Incorrect search words")
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
		log.Debug().Msgf(fmt.Sprintf("resp %+v", resp))

		rawData, err := json.Marshal(resp)
		if err != nil {
			log.Err(err).Interface("json data", resp).Msg("error while marshalling data to json")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		if _, err = fmt.Fprint(w, string(rawData)); err != nil {
			log.Printf("error %s while writing data %s do json\n", err, string(rawData))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})

	wapp.Run()
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
	file, err := os.Open(fn)
	if err != nil {
		log.Err(err).Str("filename", fn).Msg("can't open file")
		return
	}

	defer file.Close()

	ind.MapAndCleanWords(file, fn)
	return
}

func checkFlags(c *cli.Context, str ...string) error {
	for _, flag := range str {
		if c.String(flag) == "" {
			return errors.New(fmt.Sprintf("empty flag %s", flag))
		}
	}
	return nil
}
