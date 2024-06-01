package utils

import (
	"testing"
	"time"
)

func TestTTLList_Add(t *testing.T) {
	list := NewTTLList()

	// Add an item with a TTL of 1 second
	list.Add("item1", time.Second)

	// Check if the item exists in the list
	if !list.Contains("item1") {
		t.Errorf("Expected item1 to exist in the list")
	}

	// Wait for the item to expire
	time.Sleep(time.Second)

	// Check if the item has been removed from the list
	if list.Contains("item1") {
		t.Errorf("Expected item1 to be removed from the list")
	}
}

func TestTTLList_Remove(t *testing.T) {
	list := NewTTLList()

	// Add an item to the list
	list.Add("item1", time.Minute)

	// Check if the item exists in the list
	if !list.Contains("item1") {
		t.Errorf("Expected item1 to exist in the list")
	}

	// Remove the item from the list
	list.Remove("item1")

	// Check if the item has been removed
	if list.Contains("item1") {
		t.Errorf("Expected item1 to be removed from the list")
	}
}

func TestTTLList_Contains(t *testing.T) {
	list := NewTTLList()

	// Add an item to the list
	list.Add("item1", time.Minute)

	// Check if the item exists in the list
	if !list.Contains("item1") {
		t.Errorf("Expected item1 to exist in the list")
	}

	// Check if a non-existing item exists in the list
	if list.Contains("item2") {
		t.Errorf("Expected item2 to not exist in the list")
	}
}

func TestTTLList_Extend(t *testing.T) {
	list := NewTTLList()

	// Add an item with a TTL of 5 second
	list.Add("item1", time.Second)

	// Check if the item exists in the list
	if !list.Contains("item1") {
		t.Errorf("Expected item1 to exist in the list")
	}
	// Wait for the item to almost expire
	time.Sleep(50 * time.Millisecond)

	// Extend the TTL of the item with 2 seconds
	list.Add("item1", time.Second)
	// Wait for the item to expire
	time.Sleep(time.Second)
	//at this time first timer should have expired and second timer should be running
	if !list.Contains("item1") {
		t.Errorf("Expected item1 to exist in the list")
	}
	// Wait for the item to expire
	time.Sleep(time.Second)

	// Check if the item has been removed from the list
	if list.Contains("item1") {
		t.Errorf("Expected item1 to be removed from the list")
	}
}
