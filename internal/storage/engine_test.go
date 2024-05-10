package storage

import (
	"context"
	"testing"
)

func TestEngine_Set(t *testing.T) {
	storage := New()
	res, err := storage.Set(context.Background(), "testKey", "testValue")

	if err != nil {
		t.Fatalf("Set returned an error: %s", err)
	}

	if res != "ok" {
		t.Fatalf("Set returned unexpected result: want %s, got %s", "ok", res)
	}

	if storage.m["testKey"] != "testValue" {
		t.Fatalf("Set didn't set the correct value: want %s, got %s", "testValue", storage.m["testKey"])
	}
}

func TestEngine_Get(t *testing.T) {
	storage := New()
	storage.m["testKey"] = "testValue"

	res, err := storage.Get(context.Background(), "testKey")

	if err != nil {
		t.Fatalf("Get returned an error: %s", err)
	}

	if res != "testValue" {
		t.Fatalf("Get returned unexpected value: want %s, got %s", "testValue", res)
	}

	// Test get on a non-existent key
	res, err = storage.Get(context.Background(), "non_existent_key")

	if err != nil {
		t.Fatalf("Get returned an error: %s", err)
	}

	if res != "" {
		t.Fatalf("Get returned unexpected value: want %s, got %s", "", res)
	}
}

func TestEngine_Del(t *testing.T) {
	storage := New()
	storage.m["testKey"] = "testValue"

	res, err := storage.Del(context.Background(), "testKey")

	if err != nil {
		t.Fatalf("Del returned an error: %s", err)
	}

	if res != "ok" {
		t.Fatalf("Del returned unexpected result: want %s, got %s", "ok", res)
	}

	if _, ok := storage.m["testKey"]; ok {
		t.Fatal("Key was not deleted properly by Del function")
	}
}
