package db

import (
	"encoding/binary"
	"errors"
	"github.com/boltdb/bolt"
	"path"
)

type BoltDB struct {
	db       *bolt.DB
	filePath string
}

var BKTHeight = []byte("Height")

func NewBoltDB(dir string) (bdb *BoltDB, err error) {

	if dir == "" {
		err = errors.New("db dir is empty")
		return
	}
	filePath := path.Join(dir, "bolt.bin")
	db, err := bolt.Open(filePath, 0644, &bolt.Options{InitialMmapSize: 500000})
	if err != nil {
		return
	}

	err = db.Update(func(btx *bolt.Tx) error {
		_, err := btx.CreateBucketIfNotExists(BKTHeight)
		return err
	})
	if err != nil {
		return
	}
	bdb = &BoltDB{db: db, filePath: filePath}
	return
}

func (w *BoltDB) UpdateArbHeight(h uint64) error {

	raw := make([]byte, 8)
	binary.LittleEndian.PutUint64(raw, h)

	return w.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(BKTHeight)
		return bkt.Put([]byte("arb_height"), raw)
	})
}

func (w *BoltDB) GetArbHeight() uint64 {

	var h uint64
	_ = w.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(BKTHeight)
		raw := bkt.Get([]byte("arb_height"))
		if len(raw) == 0 {
			h = 0
			return nil
		}
		h = binary.LittleEndian.Uint64(raw)
		return nil
	})
	return h
}

func (w *BoltDB) Close() {
	w.db.Close()
}
