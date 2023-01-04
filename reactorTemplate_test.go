package inflate

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestReactorTemplate(t *testing.T) {
	invocationCount := 0

	tpl := NewTemplate(
		12345,
		func (i int) string {
			invocationCount++
			return strconv.Itoa(i)
		},
	)

	if r, err := tpl(); assert.NoError(t, err) {
		var s string
		if err := r.Get(&s); assert.NoError(t, err) {
			assert.Equal(t, "12345", s)
			assert.Equal(t, 1, invocationCount)
		}
		if err := r.Get(&s); assert.NoError(t, err) {
			assert.Equal(t, "12345", s)
			assert.Equal(t, 1, invocationCount)
		}
	}

	if r, err := tpl(); assert.NoError(t, err) {
		var s string
		if err := r.Get(&s); assert.NoError(t, err) {
			assert.Equal(t, "12345", s)
			assert.Equal(t, 2, invocationCount)
		}
		if err := r.Get(&s); assert.NoError(t, err) {
			assert.Equal(t, "12345", s)
			assert.Equal(t, 2, invocationCount)
		}
	}
}
