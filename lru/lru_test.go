package lru

import "testing"

import "reflect"

func TestSet_WithRawInteger(t *testing.T) {
	type DummyCacheValue struct {
		value string
	}

	c := New(3)

	c.Set("keyA", 7)

	want := 7
	got, ok := c.Get("keyA")
	if !ok || !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, but got %v", want, got)
	}

	_, ok = c.Get("keyBWillBeMissed")
	if ok {
		t.Errorf("want a miss but got a hit")
	}
}

func TestSet_WithSimpleStruct(t *testing.T) {
	type DummyCacheValue struct {
		value string
	}

	c := New(3)

	c.Set("keyA", DummyCacheValue{"dummyValue"})

	want := DummyCacheValue{"dummyValue"}
	got, ok := c.Get("keyA")
	if !ok || !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, but got %v", want, got)
	}

	_, ok = c.Get("keyBWillBeMissed")
	if ok {
		t.Errorf("want a miss but got a hit")
	}
}

func TestSet_WithOldestDeletedViaSetOrder(t *testing.T) {
	type DummyCacheValue struct {
		value string
	}

	c := New(2)

	c.Set("keyA", DummyCacheValue{"dummyValueA"})
	c.Set("keyB", DummyCacheValue{"dummyValueB"})
	c.Set("keyC", DummyCacheValue{"dummyValueC"})

	want := []string{"keyB", "keyC"}
	got := c.GetAllKeys()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, but got %v", want, got)
	}
}

func TestSet_WithOldestDeletedAfterAKeyMovedToFrontViaGet(t *testing.T) {
	type DummyCacheValue struct {
		value string
	}

	c := New(2)

	c.Set("keyA", DummyCacheValue{"dummyValueA"})
	c.Set("keyB", DummyCacheValue{"dummyValueB"})
	c.Set("keyC", DummyCacheValue{"dummyValueC"})

	c.Get("keyA")

	want := []string{"keyB", "keyC"}
	got := c.GetAllKeys()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, but got %v", want, got)
	}
}

func TestDelete(t *testing.T) {
	type DummyCacheValue struct {
		value string
	}

	c := New(2)

	c.Set("keyA", DummyCacheValue{"dummyValueA"})

	c.Delete("keyBMissed")
	want := DummyCacheValue{"dummyValueA"}
	got, ok := c.Get("keyA")
	if !ok || !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, but got %v", want, got)
	}

	c.Delete("keyA")
	_, ok = c.Get("keyA")
	if ok {
		t.Errorf("want a miss but got a hit")
	}
}
