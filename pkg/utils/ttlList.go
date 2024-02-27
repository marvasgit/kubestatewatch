package utils

import (
	"strings"
	"sync"
	"time"
)

// ItemSlice is a type alias for a slice of Items
type ItemSlice []Item

// Contains checks if a value exists in the slice.
func (items ItemSlice) ExtendIfExists(value string, ttl time.Duration) bool {
	for _, item := range items {
		if item.Value == value {
			if ttl > 0 {
				item.ExpiresAt = time.Now().Add(ttl)
			}
			return true
		}
	}
	return false
}

// TTLList represents a list with TTL for each item.
type TTLList struct {
	mu    sync.Mutex
	items ItemSlice
}

// Item represents an item in the list with a TTL.
type Item struct {
	Value     string
	ExpiresAt time.Time
}

// NewTTLList creates a new TTLList.
func NewTTLList() *TTLList {
	l := &TTLList{}
	go l.cleanupLoop()
	return l
}

// Add adds a new item to the list with a specified TTL.
func (l *TTLList) Add(value string, ttl time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	lowercaseValue := strings.ToLower(value)
	// Check if value already exists in the list
	if l.items.ExtendIfExists(lowercaseValue, ttl) {
		return
	}

	l.items = append(l.items, Item{
		Value:     lowercaseValue,
		ExpiresAt: time.Now().Add(ttl),
	})
}

// Contains checks if a value exists in the list.
func (l *TTLList) Contains(value string) bool {
	return l.items.ExtendIfExists(strings.ToLower(value), 0)
}

// Remove removes items from the list based on a matching value.
func (l *TTLList) Remove(value string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var remainingItems ItemSlice
	for _, item := range l.items {
		if item.Value != value {
			remainingItems = append(remainingItems, item)
		}
	}
	l.items = remainingItems
}

// cleanupLoop runs periodically to remove expired items.
func (l *TTLList) cleanupLoop() {
	for {
		time.Sleep(2 * time.Second) // Cleanup interval
		l.removeExpired()
	}
}

// removeExpired removes expired items from the list.
func (l *TTLList) removeExpired() {
	l.mu.Lock()
	defer l.mu.Unlock()

	currentTime := time.Now()
	var validItems []Item
	for _, item := range l.items {
		if item.ExpiresAt.After(currentTime) {
			validItems = append(validItems, item)
		}
	}
	l.items = validItems
}
