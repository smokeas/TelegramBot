package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type UserData struct {
	Todos   []string
	Notes   []string
	Finance []string
}

type Store struct {
	dir   string
	cache map[int64]*UserData
}

func NewStore(dir string) *Store {
	os.MkdirAll(dir, 0755)
	return &Store{
		dir:   dir,
		cache: make(map[int64]*UserData),
	}
}

func (s *Store) filename(userID int64) string {
	return filepath.Join(s.dir, fmt.Sprintf("%d.json", userID))
}

func (s *Store) load(userID int64) *UserData {
	if data, ok := s.cache[userID]; ok {
		return data
	}

	file := s.filename(userID)
	data := &UserData{}
	bytes, err := os.ReadFile(file)
	if err == nil {
		_ = json.Unmarshal(bytes, data)
	}
	s.cache[userID] = data
	return data
}

func (s *Store) save(userID int64) {
	file := s.filename(userID)
	data, _ := json.MarshalIndent(s.cache[userID], "", "  ")
	_ = os.WriteFile(file, data, 0644)
}

// === TODO ===
func (s *Store) AddTodo(userID int64, text string) {
	u := s.load(userID)
	u.Todos = append(u.Todos, text)
	s.save(userID)
}

func (s *Store) ListTodos(userID int64) string {
	u := s.load(userID)
	if len(u.Todos) == 0 {
		return ""
	}
	var out strings.Builder
	for i, t := range u.Todos {
		fmt.Fprintf(&out, "%d. %s\n", i+1, t)
	}
	return out.String()
}

func (s *Store) DoneTodo(userID int64, indexStr string) {
	u := s.load(userID)
	i, err := strconv.Atoi(indexStr)
	if err != nil || i < 1 || i > len(u.Todos) {
		return
	}
	u.Todos[i-1] = "✅ " + u.Todos[i-1]
	s.save(userID)
}

func (s *Store) DeleteTodo(userID int64, indexStr string) {
	u := s.load(userID)
	i, err := strconv.Atoi(indexStr)
	if err != nil || i < 1 || i > len(u.Todos) {
		return
	}
	u.Todos = append(u.Todos[:i-1], u.Todos[i:]...)
	s.save(userID)
}

// === NOTES ===
func (s *Store) AddNote(userID int64, text string) {
	u := s.load(userID)
	u.Notes = append(u.Notes, text)
	s.save(userID)
}

func (s *Store) ListNotes(userID int64) string {
	u := s.load(userID)
	if len(u.Notes) == 0 {
		return ""
	}
	var out strings.Builder
	for i, n := range u.Notes {
		fmt.Fprintf(&out, "%d. %s\n", i+1, n)
	}
	return out.String()
}

func (s *Store) DeleteNote(userID int64, indexStr string) {
	u := s.load(userID)
	i, err := strconv.Atoi(indexStr)
	if err != nil || i < 1 || i > len(u.Notes) {
		return
	}
	u.Notes = append(u.Notes[:i-1], u.Notes[i:]...)
	s.save(userID)
}

// === FINANCE ===
func (s *Store) AddFinance(userID int64, text string) {
	u := s.load(userID)
	u.Finance = append(u.Finance, text)
	s.save(userID)
}

func (s *Store) ListFinance(userID int64) string {
	u := s.load(userID)
	if len(u.Finance) == 0 {
		return "Нет финансовых записей."
	}
	var out strings.Builder
	for i, f := range u.Finance {
		fmt.Fprintf(&out, "%d. %s\n", i+1, f)
	}
	return out.String()
}

func (s *Store) Balance(userID int64) string {
	u := s.load(userID)
	total := 0
	for _, f := range u.Finance {
		fields := strings.Fields(f)
		if len(fields) > 0 {
			val, err := strconv.Atoi(strings.Trim(fields[0], "+-"))
			if err == nil {
				if strings.HasPrefix(fields[0], "-") {
					total -= val
				} else {
					total += val
				}
			}
		}
	}
	return fmt.Sprintf("Баланс: %d₽", total)
}

func (s *Store) Random(userID int64) string {
	u := s.load(userID)
	all := append([]string{}, u.Todos...)
	all = append(all, u.Notes...)
	if len(all) == 0 {
		return "Нет данных для выбора."
	}
	rand.Seed(time.Now().UnixNano())
	return all[rand.Intn(len(all))]
}
