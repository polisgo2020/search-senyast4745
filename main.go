package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/polisgo2020/search-senyast4745/config"
	"github.com/polisgo2020/search-senyast4745/database"
	"github.com/polisgo2020/search-senyast4745/web"
	"github.com/urfave/cli/v2"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/xlab/closer"

	"github.com/polisgo2020/search-senyast4745/index"
)

func main() {

	defer closer.Close()

	var err error

	if err = initLogger(config.Load()); err != nil {
		log.Err(err).Msg("can not init logger")
		return
	}

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
		DefaultText: "",
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
		log.Err(err).Msg("fatal while starting command line app")
	}
}

func initLogger(c *config.Config) error {
	logLvl, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		return err
	}
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(logLvl)
	return nil
}

func build(c *cli.Context) error {

	log.Info().Msg("build mode run")

	log.Debug().
		Str("index file", c.String("index")).
		Str("source folder", c.String("sources")).
		Msg("build run")

	if err := checkFlags(c, "sources"); err != nil {
		log.Err(err).Strs("context flags", c.FlagNames()).Msg("error while checking context")
		return nil
	}
	if allFiles, err := filePathWalkDir(c.String("sources")); err != nil {
		log.Err(err).Str(" directory", c.String("sources")).Msg("can not read files list")
	} else {
		log.Debug().Strs("files", allFiles).Msg("folder parsed")

		m := collectWordData(allFiles)

		log.Debug().Interface("index", m).Msg("index built")

		if c.String("index") != "" {
			if err := collectAndWriteMap(m, c.String("index")); err != nil {
				log.Err(err).Str("filename", c.String("index")).Msg("can not save data to file")
				return nil
			}
			log.Info().Msg("index saved")

		} else {
			repo, err := database.NewIndexRepository(config.Load())
			if err != nil {
				log.Err(err).Msg("can not open database connection")
				return nil
			}

			if err := repo.DropIndex(); err != nil {
				log.Err(err).Msg("can not drop index connection")
				return nil
			}

			if err := repo.SaveIndex(m); err != nil {
				log.Err(err).Msg("can not save index")
				return nil
			}
		}
	}

	log.Info().Msg("build done")

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

type DbIndexed struct {
	repo *database.IndexRepository
}

func (d *DbIndexed) GetIndex(str ...string) (*index.Index, error) {
	return d.repo.FindAllByWords(str)
}

type FileIndexed struct {
	i *index.Index
}

func (f *FileIndexed) GetIndex(_ ...string) (*index.Index, error) {
	return f.i, nil
}

func search(c *cli.Context) error {

	log.Info().Str("test", "Hello world").Msg("search mode run")

	log.Debug().Str("index file", c.String("index")).Msg("search run")

	cfg := config.Load()

	var wapp *web.App
	if c.String("index") != "" {
		data, err := readCSVFile(c.String("index"))
		if err != nil {
			log.Err(err).Str("file", c.String("index")).Msg("couldn't open or read the csv file")
			return nil
		}

		wapp, err = web.NewApp(cfg, &FileIndexed{i: data})
		if err != nil {
			log.Err(err).Msg("couldn't start web app")
			return nil
		}
	} else {
		repo, err := database.NewIndexRepository(cfg)
		if err != nil {
			log.Err(err).Msg("can not open database connection")
			return nil
		}
		wapp, err = web.NewApp(cfg, &DbIndexed{repo: repo})

		if err != nil {
			log.Err(err).Msg("error while creating web application")
			return nil
		}
	}
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
}

func checkFlags(c *cli.Context, str ...string) error {
	for _, flag := range str {
		if c.String(flag) == "" {
			return fmt.Errorf("empty flag %s", flag)
		}
	}
	return nil
}
