package main

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Store struct {
	dir string
	mu  sync.Mutex
}

type UserData struct {
	Pages   []string      `json:"pages"`
	Todos   []Todo        `json:"todos"`
	Notes   []string      `json:"notes"`
	Finance []Transaction `json:"finance"`
}

type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type Transaction struct {
	Kind   string  `json:"kind"` // income | expense
	Amount float64 `json:"amount"`
	Desc   string  `json:"desc"`
}

func NewStore(dir string) *Store {
	os.MkdirAll(dir, 0o755)
	return &Store{dir: dir}
}

func (s *Store) file(userID int64) string {
	return filepath.Join(s.dir, strconv.FormatInt(userID, 10)+".json")
}

func (s *Store) load(userID int64) (UserData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.file(userID)
	var data UserData

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return data, nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return data, err
	}

	return data, nil
}

func (s *Store) save(userID int64, data UserData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.file(userID)
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0o644)
}

// --- Pages ---
func (s *Store) AddPage(userID int64, url string) error {
	data, _ := s.load(userID)
	data.Pages = append(data.Pages, url)
	return s.save(userID, data)
}

func (s *Store) PickRandomPage(userID int64) (string, error) {
	data, _ := s.load(userID)
	if len(data.Pages) == 0 {
		return "", errors.New("no pages")
	}
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(data.Pages))
	page := data.Pages[idx]
	data.Pages = append(data.Pages[:idx], data.Pages[idx+1:]...)
	s.save(userID, data)
	return page, nil
}

// --- Todos ---
func (s *Store) AddTodo(userID int64, task string) {
	data, _ := s.load(userID)
	data.Todos = append(data.Todos, Todo{Text: task})
	s.save(userID, data)
}

func (s *Store) GetTodos(userID int64) []Todo {
	data, _ := s.load(userID)
	return data.Todos
}

func (s *Store) MarkTodoDone(userID int64, idx int) error {
	data, _ := s.load(userID)
	if idx < 0 || idx >= len(data.Todos) {
		return errors.New("out of range")
	}
	data.Todos[idx].Done = true
	return s.save(userID, data)
}

func (s *Store) DeleteTodo(userID int64, idx int) error {
	data, _ := s.load(userID)
	if idx < 0 || idx >= len(data.Todos) {
		return errors.New("out of range")
	}
	data.Todos = append(data.Todos[:idx], data.Todos[idx+1:]...)
	return s.save(userID, data)
}

// --- Notes ---
func (s *Store) AddNote(userID int64, text string) {
	data, _ := s.load(userID)
	data.Notes = append(data.Notes, text)
	s.save(userID, data)
}

func (s *Store) GetNotes(userID int64) []string {
	data, _ := s.load(userID)
	return data.Notes
}

func (s *Store) DeleteNote(userID int64, idx int) error {
	data, _ := s.load(userID)
	if idx < 0 || idx >= len(data.Notes) {
		return errors.New("out of range")
	}
	data.Notes = append(data.Notes[:idx], data.Notes[idx+1:]...)
	return s.save(userID, data)
}

// --- Finance ---
func (s *Store) AddFinance(userID int64, kind string, amount float64, desc string) error {
	if kind != "income" && kind != "expense" {
		return errors.New("invalid kind")
	}
	data, _ := s.load(userID)
	data.Finance = append(data.Finance, Transaction{Kind: kind, Amount: amount, Desc: desc})
	return s.save(userID, data)
}

func (s *Store) GetFinance(userID int64) []Transaction {
	data, _ := s.load(userID)
	return data.Finance
}

func (s *Store) FinanceBalance(userID int64) (income, expense float64) {
	data, _ := s.load(userID)
	for _, t := range data.Finance {
		if t.Kind == "income" {
			income += t.Amount
		} else {
			expense += t.Amount
		}
	}
	return
}
