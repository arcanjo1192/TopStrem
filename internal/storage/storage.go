package storage

import (
    "encoding/json"
    "errors"
    "time"
    "strings"

    bolt "go.etcd.io/bbolt"
)

var (
    usersBucket        = []byte("users")
    favoritesBucket    = []byte("favorites")
    watchedBucket      = []byte("watched_episodes")
	listsBucket 	   = []byte("lists")
)

type UserProfile struct {
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    LastLogin time.Time `json:"lastLogin"`
}

type FavoriteItem struct {
    ID    string `json:"id"`
    Type  string `json:"type"`
    Name  string `json:"name"`
    Year  string `json:"year"`
}

type ListInfo struct {
    Name  string         `json:"name"`
    Type  string         `json:"type"`
    Items []FavoriteItem `json:"items"`
}

type Storage struct {
    db *bolt.DB
}

func Open(path string) (*Storage, error) {
    db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: 1 * time.Second})
    if err != nil {
        return nil, err
    }

    err = db.Update(func(tx *bolt.Tx) error {
        if _, err := tx.CreateBucketIfNotExists(usersBucket); err != nil {
            return err
        }
        if _, err := tx.CreateBucketIfNotExists(favoritesBucket); err != nil {
            return err
        }
        if _, err := tx.CreateBucketIfNotExists(watchedBucket); err != nil {
            return err
        }
		if _, err := tx.CreateBucketIfNotExists(listsBucket); err != nil {
			return err
		}
        return nil
    })
    if err != nil {
        db.Close()
        return nil, err
    }

    return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
    if s.db == nil {
        return nil
    }
    return s.db.Close()
}

func (s *Storage) SaveUser(profile UserProfile) error {
    if profile.Email == "" {
        return errors.New("email is required")
    }
    profile.LastLogin = time.Now().UTC()
    return s.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(usersBucket)
        data, err := json.Marshal(profile)
        if err != nil {
            return err
        }
        return bucket.Put([]byte(profile.Email), data)
    })
}

func (s *Storage) GetUser(email string) (*UserProfile, error) {
    if email == "" {
        return nil, errors.New("email is required")
    }
    var profile UserProfile
    err := s.db.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(usersBucket)
        data := bucket.Get([]byte(email))
        if data == nil {
            return errors.New("user not found")
        }
        return json.Unmarshal(data, &profile)
    })
    if err != nil {
        return nil, err
    }
    return &profile, nil
}

func (s *Storage) SaveFavorites(email string, items []FavoriteItem) error {
    if email == "" {
        return errors.New("email is required")
    }
    return s.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(favoritesBucket)
        data, err := json.Marshal(items)
        if err != nil {
            return err
        }
        return bucket.Put([]byte(email), data)
    })
}

func (s *Storage) GetFavorites(email string) ([]FavoriteItem, error) {
    if email == "" {
        return nil, errors.New("email is required")
    }

    var items []FavoriteItem
    err := s.db.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(favoritesBucket)
        data := bucket.Get([]byte(email))
        if data == nil {
            items = []FavoriteItem{}
            return nil
        }
        return json.Unmarshal(data, &items)
    })
    if err != nil {
        return nil, err
    }
    return items, nil
}

func (s *Storage) AddFavorite(email string, item FavoriteItem) error {
    if email == "" || item.ID == "" || item.Type == "" {
        return errors.New("email, id and type are required")
    }
    existing, err := s.GetFavorites(email)
    if err != nil {
        return err
    }
    for _, current := range existing {
        if current.ID == item.ID {
            return nil
        }
    }
    existing = append(existing, item)
    return s.SaveFavorites(email, existing)
}

func (s *Storage) RemoveFavorite(email string, id string) error {
    if email == "" || id == "" {
        return errors.New("email and id are required")
    }
    existing, err := s.GetFavorites(email)
    if err != nil {
        return err
    }
    filtered := make([]FavoriteItem, 0, len(existing))
    for _, current := range existing {
        if current.ID != id {
            filtered = append(filtered, current)
        }
    }
    return s.SaveFavorites(email, filtered)
}

func (s *Storage) SaveWatchedEpisodes(email string, episodeIDs []string) error {
    if email == "" {
        return errors.New("email is required")
    }
    return s.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(watchedBucket)
        data, err := json.Marshal(episodeIDs)
        if err != nil {
            return err
        }
        return bucket.Put([]byte(email), data)
    })
}

func (s *Storage) GetWatchedEpisodes(email string) ([]string, error) {
    if email == "" {
        return nil, errors.New("email is required")
    }

    var episodes []string
    err := s.db.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(watchedBucket)
        data := bucket.Get([]byte(email))
        if data == nil {
            episodes = []string{}
            return nil
        }
        return json.Unmarshal(data, &episodes)
    })
    if err != nil {
        return nil, err
    }
    return episodes, nil
}

func (s *Storage) AddWatchedEpisode(email string, episodeID string) error {
    if email == "" || episodeID == "" {
        return errors.New("email and episodeID are required")
    }
    existing, err := s.GetWatchedEpisodes(email)
    if err != nil {
        return err
    }
    for _, current := range existing {
        if current == episodeID {
            return nil
        }
    }
    existing = append(existing, episodeID)
    return s.SaveWatchedEpisodes(email, existing)
}

func (s *Storage) RemoveWatchedEpisode(email string, episodeID string) error {
    if email == "" || episodeID == "" {
        return errors.New("email and episodeID are required")
    }
    existing, err := s.GetWatchedEpisodes(email)
    if err != nil {
        return err
    }
    filtered := make([]string, 0, len(existing))
    for _, current := range existing {
        if current != episodeID {
            filtered = append(filtered, current)
        }
    }
    return s.SaveWatchedEpisodes(email, filtered)
}

func (s *Storage) GetList(email, listName string) (*ListInfo, error) {
    if email == "" || listName == "" {
        return nil, errors.New("email and list name are required")
    }
    key := email + ":" + listName
    var list ListInfo
    err := s.db.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(listsBucket)
        data := bucket.Get([]byte(key))
        if data == nil {
            return errors.New("list not found")
        }
        return json.Unmarshal(data, &list)
    })
    if err != nil {
        return nil, err
    }
    return &list, nil
}

func (s *Storage) GetAllLists(email string) ([]ListInfo, error) {
    if email == "" {
        return nil, errors.New("email is required")
    }
    var lists []ListInfo
    err := s.db.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(listsBucket)
        cursor := bucket.Cursor()
        prefix := []byte(email + ":")
        for k, v := cursor.Seek(prefix); k != nil && strings.HasPrefix(string(k), string(prefix)); k, v = cursor.Next() {
            var list ListInfo
            if err := json.Unmarshal(v, &list); err != nil {
                return err
            }
            lists = append(lists, list)
        }
        return nil
    })
    return lists, err
}

func (s *Storage) CreateList(email, listName, listType string) error {
    if email == "" || listName == "" || listType == "" {
        return errors.New("email, list name and type are required")
    }
    key := email + ":" + listName
    return s.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(listsBucket)
        if bucket.Get([]byte(key)) != nil {
            return errors.New("list already exists")
        }
        list := ListInfo{
            Name:  listName,
            Type:  listType,
            Items: []FavoriteItem{},
        }
        data, err := json.Marshal(list)
        if err != nil {
            return err
        }
        return bucket.Put([]byte(key), data)
    })
}

func (s *Storage) DeleteList(email, listName string) error {
    if email == "" || listName == "" {
        return errors.New("email and list name are required")
    }
    key := email + ":" + listName
    return s.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(listsBucket)
        return bucket.Delete([]byte(key))
    })
}

func (s *Storage) AddItemToList(email, listName string, item FavoriteItem) error {
    if email == "" || listName == "" || item.ID == "" || item.Type == "" {
        return errors.New("invalid parameters")
    }
    key := email + ":" + listName
    return s.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(listsBucket)
        data := bucket.Get([]byte(key))
        if data == nil {
            return errors.New("list not found")
        }
        var list ListInfo
        if err := json.Unmarshal(data, &list); err != nil {
            return err
        }
        // Verifica duplicata
        for _, it := range list.Items {
            if it.ID == item.ID {
                return nil // já existe
            }
        }
        list.Items = append(list.Items, item)
        updated, err := json.Marshal(list)
        if err != nil {
            return err
        }
        return bucket.Put([]byte(key), updated)
    })
}

func (s *Storage) RemoveItemFromList(email, listName, itemID string) error {
    if email == "" || listName == "" || itemID == "" {
        return errors.New("invalid parameters")
    }
    key := email + ":" + listName
    return s.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(listsBucket)
        data := bucket.Get([]byte(key))
        if data == nil {
            return errors.New("list not found")
        }
        var list ListInfo
        if err := json.Unmarshal(data, &list); err != nil {
            return err
        }
        filtered := make([]FavoriteItem, 0, len(list.Items))
        for _, it := range list.Items {
            if it.ID != itemID {
                filtered = append(filtered, it)
            }
        }
        list.Items = filtered
        updated, err := json.Marshal(list)
        if err != nil {
            return err
        }
        return bucket.Put([]byte(key), updated)
    })
}