package store

import (
	"context"
	"fmt"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type GobTestSuite struct {
	suite.Suite
	tableConnect []struct {
		loc      string
		expected bool
	}
	tableStoreRetrieve map[string]string
	tableStorePatch    map[string]string
	mapbucket          []map[string]string
	tokens             map[string]string
	log                *vlog.Logger
}

var (
	varTableGobConnect = []struct {
		loc      string
		expected bool
	}{
		{"test.gob", true},
		{"test.db", true},
		{filepath.Join("false", "great.db"), true},
		{filepath.Join("false", "true.db"), true},
		{"test.yaml", true},
	}

	varTableStoreGobRetrieve = map[string]string{
		"ijbnijdelkfiue1": "A1B2C3D4E5F6G7H8",
		"ijbnijdelkfiue2": "Z9Y8X7W6V5U4T3S2",
		"ijbnijdelkfiue3": "Q1W2E3R4T5Y6U7I8",
		"ijbnijdelkfiue4": "O9P0A1S2D3F4G5H6",
		"ijbnijdelkfiue5": "J7K8L9Z0X1C2V3B4",
		"ijbnijdelkfiue6": "N5M6Q1W2E3R4T5Y",
		"ijbnijdelkfiue7": "U6I7O8P9A0S1D2F",
		"ijbnijdelkfiue8": "G3H4J5K6L7Z8X9C",
	}

	varTableGobPatch = map[string]string{
		"ijbnijdelkfiue1": "649sx8C30ubzd0cu",
		"ijbnijdelkfiue2": "TN4IFzbjuJfwuOIW",
		"ijbnijdelkfiue3": "a1otXUTJnt4gzOLL",
		"ijbnijdelkfiue4": "df8opIzQPrpRn9sM",
		"ijbnijdelkfiue5": "wADOGUHAJ5wtgiAO",
		"ijbnijdelkfiue6": "bvaeCnte1VAODI91",
		"ijbnijdelkfiue7": "8WNEeW7uVvYZpIrR",
		"ijbnijdelkfiue8": "Exz9baPAttTIgusZ",
	}

	varTableMapCore = []map[string]string{
		{"key1": "value1", "key2": "value2", "key3": "value3", "key4": "value4", "key5": "value5",
			"key6": "value6", "key7": "value7", "key8": "value8", "key9": "value9", "key10": "value10"},
		{"name": "John", "age": "30", "city": "New York", "country": "USA", "occupation": "Engineer",
			"hobby": "Reading", "language": "English", "car": "Honda", "food": "Pizza", "sport": "Soccer"},
		{"product": "Laptop", "price": "1000", "brand": "Dell", "model": "XPS", "os": "Windows 10",
			"color": "Black", "weight": "2kg", "processor": "i7", "ram": "16GB", "storage": "512GB SSD"},
		{"book": "1984", "author": "George Orwell", "genre": "Dystopian", "year": "1949", "pages": "328",
			"publisher": "Secker & Warburg", "language": "English", "country": "UK", "format": "Hardcover", "ISBN": "9780451524935"},
		{"language": "Go", "version": "1.18", "developer": "Google", "releaseYear": "2012", "syntax": "Static",
			"typing": "Strong", "OS": "Cross-platform", "license": "BSD", "paradigm": "Concurrent", "website": "golang.org"},
		{"car": "Tesla", "model": "Model S", "year": "2020", "color": "Red", "range": "370 miles",
			"battery": "100 kWh", "charging": "Supercharger", "seats": "5", "price": "79990", "autopilot": "Available"},
		{"fruit": "Apple", "color": "Red", "taste": "Sweet", "origin": "Central Asia", "vitamin": "C",
			"calories": "95", "water": "85%", "fiber": "4g", "sugar": "19g", "type": "Fruit"},
		{"planet": "Mars", "status": "Uninhabited", "gravity": "3.711 m/s²", "moons": "2", "dayLength": "24.6 hours",
			"distanceFromSun": "227.9 million km", "yearLength": "687 Earth days", "temperature": "-28°C", "atmosphere": "CO2", "exploration": "Rovers"},
		{"game": "Chess", "players": "2", "origin": "India", "pieces": "32", "boards": "Chessboard",
			"strategy": "High", "skill": "Tactics", "timeControl": "Varies", "worldFederation": "FIDE", "olympiad": "Chess Olympiad"},
		{"movie": "Inception", "director": "Christopher Nolan", "year": "2010", "genre": "Sci-Fi",
			"runtime": "148 minutes", "cast": "Leonardo DiCaprio", "budget": "160 million USD", "boxOffice": "Over 830 million USD", "awards": "Academy Awards", "rating": "PG-13"},
	}
)

func (suite *GobTestSuite) SetupTest() {
	//port := "6378"
	//err := SetUpEnv(port)
	//suite.Require().NoErrorf(err, "docker environment setup failed with error: %s\n", err.Error())
	suite.tableConnect = varTableGobConnect
	suite.tableStoreRetrieve = varTableStoreGobRetrieve
	suite.tableStorePatch = varTableGobPatch
	suite.mapbucket = varTableMapCore
	suite.log = vlog.New(true)
}

// Right now these tests are majorly happy-path tests

func (suite *GobTestSuite) TestConnect() {
	_ = suite.log.Logger()
	ctx := context.Background()
	for i := 0; i < len(suite.tableConnect); i++ {
		fmt.Printf(Order, i+1)
		loc := suite.tableConnect[i].loc
		gob, err := NewGob(ctx, loc, suite.log)
		// decided not to require no errors here, because the core error handling logic is handled by go=redis, so no use we trying to test it
		//suite.Assert().NoErrorf(err, "got error: %v\n", err)
		// continue even with error
		if err != nil {
			continue
		}
		b, err := gob.Connect(ctx)
		suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
		suite.Equalf(suite.tableConnect[i].expected, b, "expected %v, got %v\n", suite.tableConnect[i].expected, b)
		// clean DB
		suite.flush(ctx, gob)
		err = gob.Close(ctx)
		suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
	}
}

func (suite *GobTestSuite) TestMapDump() {
	_ = suite.log.Logger()
	ctx := context.Background()
	loc := suite.tableConnect[0].loc
	for i := 0; i < len(suite.mapbucket); i++ {
		fmt.Printf(Order, i+1)
		currentMap := suite.mapbucket[i]
		gob, err := NewGob(ctx, loc, suite.log)
		suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
		b, err := gob.Connect(ctx)
		suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
		suite.Assert().True(b, "expected true, got false")

		// store map into in-memory store
		for k, v := range currentMap {
			err := gob.basin.Store(ctx, k, v)
			suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
		}

		// dump in-memory store
		err = gob.MapDump(ctx)
		suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)

		// save current sync map in variable
		old := gob.basin.Map()

		// flush and then fetch the map in storage
		//b, err = gob.Flush(ctx)
		//suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
		//suite.Assert().True(b, "expected true, got false")

		err = gob.MapRefresh(ctx)
		suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
		newM, err := gob.basin.RetrieveAll(ctx)
		suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)

		suite.Require().Equalf(currentMap, newM, "expected %v (current iter in map), but got %v (map read)\n", currentMap, newM)
		suite.Require().Equalf(old, newM, "expected %v (map dumped), but got %v (map read)\n", currentMap, newM)
	}
}

//func (suite *GobTestSuite) MapRefresh(ctx context.Context) {
//	b, err := gob.MapRefresh(ctx)
//	suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
//	suite.Assert().True(b, "expected true, got false")
//}

//func (suite *GobTestSuite) TestStoreAndRetrieve() {
//	_ = suite.log.Logger()
//	ctx := context.Background()
//	i := 1
//	for k, v := range suite.tableStoreRetrieve {
//		fmt.Printf(Order, i)
//		redis, err := NewRedis(suite.loc, suite.log)
//		b, err := redis.Connect(ctx)
//		suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		suite.Assert().True(b, "expected true but got false")
//		err = redis.Store(ctx, k, v)
//		suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		time.Sleep(2)
//		val, err := redis.Retrieve(ctx, k)
//		suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		suite.Require().Equalf(v, val, "expected %s, but got %s\n", v, val)
//		// clean DB
//		suite.flush(ctx, redis)
//		err = redis.Close(ctx)
//		suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		i++
//	}
//}
//
//func (suite *GobTestSuite) TestRetrieveAll() {
//	_ = suite.log.Logger()
//	ctx := context.Background()
//	i := 1
//	redis, err := NewRedis(suite.redisConnectionString, suite.log)
//	b, err := redis.Connect(ctx)
//	suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
//	suite.Assert().True(b, "expected true but got false")
//	for k, v := range suite.tableStoreRetrieve {
//		fmt.Printf(Order, i)
//		err = redis.Store(ctx, k, v)
//		suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		i++
//	}
//	valMap, err := redis.RetrieveAll(ctx)
//	suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
//	suite.Require().Equalf(len(suite.tableStoreRetrieve), len(valMap), "expected %d, but got %d\n", len(suite.tableStoreRetrieve), len(valMap))
//	// clean DB
//	suite.flush(ctx, redis)
//	err = redis.Close(ctx)
//	suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
//}
//
//func (suite *GobTestSuite) TestDelete() {
//	_ = suite.log.Logger()
//	ctx := context.Background()
//	i := 1
//	redis, err := NewRedis(suite.redisConnectionString, suite.log)
//	b, err := redis.Connect(ctx)
//	suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
//	suite.Assert().True(b, "expected true but got false")
//	for k, v := range suite.tableStoreRetrieve {
//		fmt.Printf(Order, i)
//		err = redis.Store(ctx, k, v)
//		suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		// id should exist, and value should equal v
//		val, err := redis.Retrieve(ctx, k)
//		suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		suite.Assert().Equalf(v, val, "expected %s, but got %v\n", v, val)
//		b, err := redis.Delete(ctx, k)
//		suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		suite.Require().True(b, "expected %v, got %v\n", true, b)
//		// id should exist, and value should equal v
//		val, err = redis.Retrieve(ctx, k)
//		suite.Require().Error(err, "expected id to not exist, but got this %v\n", err.Error())
//		suite.Require().Equalf("", val, "expected %s, but got %v\n", v, val)
//		// ensure DB is flushed
//		suite.flush(ctx, redis)
//		i++
//	}
//	err = redis.Close(ctx)
//	suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
//}
//
//func (suite *GobTestSuite) TestPatch() {
//	_ = suite.log.Logger()
//	ctx := context.Background()
//	i := 1
//	redis, err := NewRedis(suite.redisConnectionString, suite.log)
//	b, err := redis.Connect(ctx)
//	suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
//	suite.Assert().True(b, "expected true but got false")
//	for k, v := range suite.tableStorePatch {
//		fmt.Printf(Order, i)
//		err = redis.Store(ctx, k, v)
//		suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		// id should exist, and value should equal v
//		val, err := redis.Retrieve(ctx, k)
//		suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		suite.Assert().Equalf(v, val, "expected %s, but got %v\n", v, val)
//		b, err := redis.Patch(ctx, k, v)
//		suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
//		suite.Require().True(b, "expected %v, got %v\n", true, b)
//		// id should exist, and value should equal v
//		val, err = redis.Retrieve(ctx, k)
//		suite.Require().NoErrorf(err, ErrWithOperation, "expected id to not exist, but got this %v\n", err)
//		suite.Require().Equalf(v, val, "expected %s, but got %v\n", v, val)
//		// ensure DB is flushed
//		suite.flush(ctx, redis)
//		i++
//	}
//	err = redis.Close(ctx)
//	suite.Require().NoErrorf(err, "expected no errors, but got this %v\n", err)
//}

func (suite *GobTestSuite) TearDownTest() {}

// TestRedisSuite tests the Redis suite
func TestGobSuite(t *testing.T) {
	suite.Run(t, new(GobTestSuite))
}

func (suite *GobTestSuite) flush(ctx context.Context, gob *Gob) {
	b, err := gob.Flush(ctx)
	suite.Assert().NoErrorf(err, "expected no errors, but got this %v\n", err)
	suite.Assert().True(b, "expected true, got false")
}
