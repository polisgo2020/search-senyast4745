package database

import (
	"context"
	"sync"
	"time"

	"github.com/polisgo2020/search-senyast4745/config"
	"github.com/polisgo2020/search-senyast4745/index"
	"github.com/rs/zerolog/log"
	"github.com/xlab/closer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	client *mongo.Client
}

var (
	instance *Connection
	once     sync.Once

	database string
)

func InitDB(c *config.Config) (*Connection, error) {
	var err error
	once.Do(func() {
		database = c.Database

		var client *mongo.Client
		opt := options.Client()
		opt.SetConnectTimeout(500 * time.Millisecond)

		opt.ApplyURI(c.DbListen)
		log.Debug().Msg("start initializing database")
		client, err = mongo.NewClient(opt)
		if err != nil {
			return
		}

		err = client.Connect(context.TODO())
		if err != nil {
			return
		}

		log.Debug().Msg("Ping database")
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			return
		}
		log.Info().Msg("connected to database!")
		instance = &Connection{client: client}
		context.TODO()

		closer.Bind(instance.Close)

	})
	return instance, err
}

func (c *Connection) Close() {
	log.Debug().Msg("Start closing database connection")
	if err := c.client.Disconnect(context.TODO()); err != nil {
		log.Err(err).Msg("error while closing db connection")
	}
	log.Info().Msg("Database disconnect successfully")
}

type IndexRepository struct {
	col *mongo.Collection
}

type IndexDTO struct {
	Word    string
	FileStr []*index.FileStruct
}

type Trainer struct {
	Name string
	Age  int
	City string
}

func TransformIndex(i *index.Index) []IndexDTO {
	var dto []IndexDTO
	log.Debug().Interface("index", i).Msg("start index transfer")
	var data IndexDTO
	for k := range (*i).Data {
		data = IndexDTO{k, i.Data[k]}
		dto = append(dto, data)
	}
	log.Debug().Interface("dto", dto).Msg("dto made")
	return dto
}

func NewIndexRepository(c *config.Config) (*IndexRepository, error) {
	con, err := InitDB(c)
	if err != nil {
		return nil, err
	}
	col := con.client.Database(database).Collection("index")
	return &IndexRepository{col: col}, nil
}

func (rep *IndexRepository) SaveIndex(i *index.Index) error {
	var transfer []interface{}
	for _, v := range TransformIndex(i) {
		transfer = append(transfer, v)
	}
	log.Debug().Interface("transfer", transfer).Msg("data")
	_, err := rep.col.InsertMany(context.TODO(), transfer)
	return err
}

func (rep *IndexRepository) FindAllByWords(wordArr []string) (*index.Index, error) {
	log.Debug().Strs("words", wordArr).Msg("start find by words")
	filter := bson.M{"word": bson.M{"$in": wordArr}}
	cursor, err := rep.col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	i := index.NewIndex()
	for cursor.Next(context.TODO()) {
		var tmp IndexDTO
		err := cursor.Decode(&tmp)
		if err != nil {
			return nil, err
		}
		log.Debug().Interface("parsed cursor", tmp).Msg("cursor parsed")
		i.Data[tmp.Word] = tmp.FileStr
	}
	log.Info().Interface("index", i).Strs("words", wordArr).Msg("index get from db")
	return i, nil
}

func (rep *IndexRepository) DropIndex() error {
	return rep.col.Drop(context.TODO())
}
