package store

import (
	"encoding/json"
	"fmt"
)

type Store struct {
	bookmarks map[string]*Bookmark
	path      string
}

func NewStore(path string) (*Store, error) {
	var store Store
	store = Store{
		path:      path,
		bookmarks: make(map[string]*Bookmark),
	}
	err := store.initBookmarks()
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *Store) GetOne(name string) (*Bookmark, error) {
	bm, exists := s.bookmarks[name]
	if !exists {
		return nil, fmt.Errorf("bookmark with name %v not found", name)
	}
	return bm, nil
}

func (s *Store) GetAll() []*Bookmark {
	list := make([]*Bookmark, 0, len(s.bookmarks))
	for _, bm := range s.bookmarks {
		list = append(list, bm)
	}
	return list
}

func (s *Store) initBookmarks() error {
	data, err := s.loadFile()
	if err != nil {
		return err
	}
	var bookmarks []*Bookmark
	err = json.Unmarshal(data, &bookmarks)
	if err != nil {
		return err
	}
	for _, bm := range bookmarks {
		s.bookmarks[bm.Name] = bm
	}
	return nil
}
