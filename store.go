package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Task struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type Store struct {
	tasks map[int64][]Task
}

func NewStore() *Store {
	return &Store{tasks: make(map[int64][]Task)}
}

func (s *Store) load(userID int64) {
	if _, ok := s.tasks[userID]; ok {
		return // уже загружено
	}
	filename := fmt.Sprintf("tasks_%d.json", userID)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		s.tasks[userID] = []Task{}
		return
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		s.tasks[userID] = []Task{}
		return
	}
	s.tasks[userID] = tasks
}

func (s *Store) save(userID int64) error {
	filename := fmt.Sprintf("tasks_%d.json", userID)
	data, err := json.MarshalIndent(s.tasks[userID], "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func (s *Store) AddTask(userID int64, text string) (int, error) {
	s.load(userID)
	task := Task{Text: text, Done: false}
	s.tasks[userID] = append(s.tasks[userID], task)
	if err := s.save(userID); err != nil {
		return 0, err
	}
	return len(s.tasks[userID]), nil
}

func (s *Store) MarkDone(userID int64, index int) error {
	s.load(userID)
	tasks := s.tasks[userID]
	if index < 1 || index > len(tasks) {
		return fmt.Errorf("Неверный номер задачи")
	}
	s.tasks[userID][index-1].Done = true
	return s.save(userID)
}

func (s *Store) DeleteTask(userID int64, index int) error {
	s.load(userID)
	tasks := s.tasks[userID]
	if index < 1 || index > len(tasks) {
		return fmt.Errorf("Неверный номер задачи")
	}
	// удаляем задачу с данным индексом
	s.tasks[userID] = append(tasks[:index-1], tasks[index:]...)
	return s.save(userID)
}
