package db

import (
	"bytes"
	"encoding/gob"
	"log"
	"strconv"

	"github.com/anonyindian/logger"
	"github.com/peterbourgon/diskv/v3"
)

var client *diskv.Diskv

func Load(l *logger.Logger) {
	l = l.Create("DATABASE")
	// Simplest transform function: put all the data files into the base dir.
	flatTransform := func(s string) []string { return []string{} }

	// Initialize a new diskv store, rooted at "my-data-dir", with a 1MB cache.
	d := diskv.New(diskv.Options{
		BasePath:     "my-data-dir",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	client = d
	defer l.ChangeLevel(logger.LevelMain).Println("LOADED")
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

func get(key string, T interface{}) {
	b, _ := client.Read(key)
	gob.NewDecoder(bytes.NewBuffer(b)).Decode(T)
}

func set(key string, T interface{}) {
	setRaw(key, encode(T))
}

func getRaw(key string) string {
	return client.ReadString(key)
}

func setRaw(key string, v interface{}) {
	vb, err := GetBytes(v)
	if err != nil {
		log.Println(err.Error())
		return
	}
	client.Write(key, vb)
}

func setBool(key string, value bool) {
	setRaw(key, strconv.FormatBool(value))
}

func getBool(key string) bool {
	return parseBool(getRaw(key))
}

func encode(v interface{}) []byte {
	buf := bytes.Buffer{}
	gob.NewEncoder(&buf).Encode(v)
	return buf.Bytes()
}

func parseBool(s string) bool {
	return s == "true"
}
