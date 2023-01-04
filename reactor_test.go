package inflate

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestSimple(t *testing.T) {
	r := New()
	_ = r.Put("Fooo", 132)

	var s string
	if err := r.Get(&s); assert.NoError(t, err) {
		assert.Equal(t, "Fooo", s)
	}

	var i int
	if err := r.Get(&i); assert.NoError(t, err) {
		assert.Equal(t, 132, i)
	}
}

func TestGetReactor(t *testing.T) {
	r := New()
	var r2 Reactor
	if err := r.Get(&r2); assert.NoError(t, err) {
		assert.Same(t, r, r2)
	}
}

func TestProvider(t *testing.T) {
	r := New()
	err := r.Put(
		func() time.Duration { return time.Minute },
		func() (time.Time, error) { return time.Unix(123456789, 0), nil },
		func() *float32 {
			f := float32(10.)
			return &f
		},
		func() (*float64, error) {
			f := 33.1
			return &f, nil
		},
		func(t time.Time) string {
			return t.UTC().String()
		},
	)

	if assert.NoError(t, err) {
		var i time.Time
		if err := r.Get(&i); assert.NoError(t, err) {
			assert.Equal(t, time.Unix(123456789, 0), i)
		}
		var d time.Duration
		if err := r.Get(&d); assert.NoError(t, err) {
			assert.Equal(t, time.Minute, d)
		}
		var f32 *float32
		if err := r.Get(&f32); assert.NoError(t, err) {
			assert.Equal(t, float32(10.), *f32)
		}
		var f64 *float64
		if err := r.Get(&f64); assert.NoError(t, err) {
			assert.Equal(t, 33.1, *f64)
		}
		var s string
		if err := r.Get(&s); assert.NoError(t, err) {
			assert.Equal(t, "1973-11-29 21:33:09 +0000 UTC", s)
		}
	}
}

func TestProviderRef(t *testing.T) {
	r, err := NewWith(
		"https://google.com",
		url.Parse,
	)
	if assert.NoError(t, err) {
		var u *url.URL
		if err := r.Get(&u); assert.NoError(t, err) {
			assert.Equal(t, "google.com", u.Host)
		}
	}
}
