package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"zakirullin/dumpbot/internal/sched"
	"zakirullin/dumpbot/pkg/tg"

	"github.com/alicebob/miniredis/v2"
)

const (
	redisLastKeyboardMsgID               = "last_keyboard"
	redisReplaceWithDefaultKeyboardMsgID = "candidate_message_id"
	redisSchedule                        = "schedule"
	redisInputExpectation                = "input_expectation"
)

// DB Maybe user ID here?
type DB struct {
	redis *miniredis.Miniredis
}

func NewDB(redis *miniredis.Miniredis) *DB {
	return &DB{redis}
}

func (db *DB) LastKeyboardMsgID(userID int64) (*int, error) {
	val, err := db.redis.Get(db.key(userID, redisLastKeyboardMsgID))
	if errors.Is(err, miniredis.ErrKeyNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("can't get last keyboard msg ID: %w", err)
	}

	i, err := strconv.Atoi(val)

	return &i, err
}

func (db *DB) SetLastKeyboardMsgID(userID int64, ID int) error {
	return db.redis.Set(db.key(userID, redisLastKeyboardMsgID), strconv.Itoa(ID))
}

func (db *DB) DelLastKeyboardMsgID(userID int64) error {
	db.redis.Del(db.key(userID, redisLastKeyboardMsgID))
	//if err {
	//	return fmt.Errorf("db.DelLastKeyboardMsgID: %w", err)
	//}

	return nil
}

// map[filename] => Cron
func (db *DB) Schedule(userID int64) (map[string]sched.Cron, error) {
	sc, err := db.redis.Get(db.key(userID, redisSchedule))
	if errors.Is(err, miniredis.ErrKeyNotFound) {
		return make(map[string]sched.Cron), nil
	} else if err != nil {
		return nil, fmt.Errorf("can't get schedule: %w", err)
	}

	var schedules map[string]sched.Cron
	err = json.Unmarshal([]byte(sc), &schedules)
	if err != nil {
		return nil, fmt.Errorf("getSchedule can't unmarshal: %w", err)
	}

	return schedules, nil
}

// Schedule task from _trash_ or later at runAt (Unix timestamp, sec). Tasks appear in today folder.
// If cron is provided this task will be repeated. Other wise, it will be executed once.
func (db *DB) AddToSchedule(userID int64, filename string, runAt int64, cron string) error {
	sc, err := db.Schedule(userID)
	if err != nil {
		return fmt.Errorf("addToSchedule: can't add to schedule: %w", err)
	}

	sc[filename] = sched.NewCron(runAt, cron)

	js, err := json.Marshal(sc)
	if err != nil {
		return fmt.Errorf("can't marshal to json: %w", err)
	}

	err = db.redis.Set(db.key(userID, redisSchedule), string(js))
	if err != nil {
		return fmt.Errorf("addToSchedule: can't save schedule: %w", err)
	}

	return nil
}

func (db *DB) DelFromSchedule(userID int64, filename string) error {
	sc, err := db.Schedule(userID)
	if err != nil {
		return fmt.Errorf("delFromSchedule: can't add to schedule: %w", err)
	}

	delete(sc, filename)

	js, err := json.Marshal(sc)
	if err != nil {
		return fmt.Errorf("delFromSchedule: can't marshal to json: %w", err)
	}

	err = db.redis.Set(db.key(userID, redisSchedule), string(js))
	if err != nil {
		return fmt.Errorf("delFromSchedule: can't save schedule: %w", err)
	}

	return nil
}

func (db *DB) InputExpectation(userID int64) (*tg.Cmd, error) {
	js, err := db.redis.Get(db.key(userID, redisInputExpectation))
	if errors.Is(err, miniredis.ErrKeyNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("db.GetInputExpectation: can't get input expectation: %w", err)
	}

	cmd := new(tg.Cmd)
	err = json.Unmarshal([]byte(js), &cmd)
	if err != nil {
		return nil, fmt.Errorf("db.GetInputExpectation: can't unmarshall cmd: %w", err)
	}

	return cmd, nil
}

func (db *DB) SetInputExpectation(userID int64, cmd tg.Cmd) error {
	js, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("db.SetInputExpectation: can't json: %w", err)
	}

	err = db.redis.Set(db.key(userID, redisInputExpectation), string(js))
	if err != nil {
		return fmt.Errorf("db.SetInputExpectation: can't set to redis: %w", err)
	}

	return nil
}

func (db *DB) DelInputExpectation(userID int64) error {
	key := db.key(userID, redisInputExpectation)
	ok := db.redis.Del(key)
	if !ok {
		return errors.New(fmt.Sprintf("db.DelInputExpecation: can't del key %s", key))
	}

	return nil
}

// User-namespaced redis key
func (db *DB) key(userID int64, key string) string {
	return fmt.Sprintf("%s:%d", key, userID)
}
