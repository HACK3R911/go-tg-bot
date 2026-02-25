package repository

import (
	"sync"
	"testing"
)

func TestAuthDB_AuthorizeRepo(t *testing.T) {
	db := NewAuthDB()

	tests := []struct {
		name     string
		userID   int64
		wantAuth bool
	}{
		{"authorize user 1", 1, true},
		{"authorize user 2", 2, true},
		{"authorize negative user", -100, true},
		{"authorize zero user", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.AuthorizeRepo(tt.userID)
			got := db.IsAuthorizedRepo(tt.userID)
			if got != tt.wantAuth {
				t.Errorf("AuthorizeRepo() = %v, want %v", got, tt.wantAuth)
			}
		})
	}
}

func TestAuthDB_IsAuthorizedRepo(t *testing.T) {
	db := NewAuthDB()

	db.AuthorizeRepo(100)
	db.AuthorizeRepo(200)

	tests := []struct {
		name     string
		userID   int64
		wantAuth bool
	}{
		{"authorized user 100", 100, true},
		{"authorized user 200", 200, true},
		{"non-authorized user 0", 0, false},
		{"non-authorized user 1", 1, false},
		{"non-authorized user 999", 999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := db.IsAuthorizedRepo(tt.userID)
			if got != tt.wantAuth {
				t.Errorf("IsAuthorizedRepo() = %v, want %v", got, tt.wantAuth)
			}
		})
	}
}

func TestAuthDB_AuthorizeRepo_MultipleTimes(t *testing.T) {
	db := NewAuthDB()
	userID := int64(12345)

	db.AuthorizeRepo(userID)
	db.AuthorizeRepo(userID)
	db.AuthorizeRepo(userID)

	got := db.IsAuthorizedRepo(userID)
	if !got {
		t.Errorf("IsAuthorizedRepo() after multiple AuthorizeRepo = false, want true")
	}
}

func TestAuthDB_ConcurrentAuthorize(t *testing.T) {
	db := NewAuthDB()
	userID := int64(999)
	iterations := 1000

	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db.AuthorizeRepo(userID)
		}()
	}
	wg.Wait()

	if !db.IsAuthorizedRepo(userID) {
		t.Error("user should be authorized after concurrent writes")
	}
}

func TestAuthDB_ConcurrentReadWrite(t *testing.T) {
	db := NewAuthDB()
	iterations := 1000

	var wg sync.WaitGroup

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			db.AuthorizeRepo(id)
		}(int64(i))
	}

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			db.IsAuthorizedRepo(id)
		}(int64(i))
	}

	wg.Wait()
}

func BenchmarkAuthDB_AuthorizeRepo(b *testing.B) {
	db := NewAuthDB()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.AuthorizeRepo(int64(i))
	}
}

func BenchmarkAuthDB_IsAuthorizedRepo(b *testing.B) {
	db := NewAuthDB()
	for i := 0; i < 1000; i++ {
		db.AuthorizeRepo(int64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.IsAuthorizedRepo(int64(i % 1000))
	}
}

func BenchmarkAuthDB_ConcurrentReadWrite(b *testing.B) {
	db := NewAuthDB()
	for i := 0; i < 1000; i++ {
		db.AuthorizeRepo(int64(i))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := 0
		for pb.Next() {
			db.AuthorizeRepo(int64(id))
			db.IsAuthorizedRepo(int64(id % 1000))
			id++
		}
	})
}
