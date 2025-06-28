package faker_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shoraid/stx-go-utils/faker"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name       string     `faker:"string"`
	Bio        *string    `faker:"sentence"`
	IsActive   bool       `faker:"bool"`
	Age        int        `faker:"int"`
	JoinedAt   *time.Time `faker:"time"`
	RelationID string     `faker:"uuid_str"`
	NoTag      string     // should remain zero-value
}

type MinimalStruct struct {
	Flag *bool `faker:"bool"`
}

func TestFaker_GenerateFake(t *testing.T) {
	tests := []struct {
		name       string
		generate   func() any
		assertions func(t *testing.T, result any)
	}{
		{
			name: "TestStruct full field",
			generate: func() any {
				return faker.GenerateFake[TestStruct]()
			},
			assertions: func(t *testing.T, result any) {
				ts := result.(*TestStruct)

				_, err := uuid.Parse(ts.RelationID)

				assert.NotEmpty(t, ts.Name, "Name should be generated")
				assert.NotNil(t, ts.Bio, "Bio should not be nil")
				assert.NotEmpty(t, *ts.Bio, "Bio should have value")
				assert.NotZero(t, ts.Age, "Age should not be zero")
				assert.NotNil(t, ts.JoinedAt, "JoinedAt should not be nil")
				assert.WithinDuration(t, time.Now(), *ts.JoinedAt, 30*24*time.Hour)
				assert.NoError(t, err, "RelationID should be valid UUID")
				assert.Empty(t, ts.NoTag, "NoTag should be empty since no faker tag")
			},
		},
		{
			name: "MinimalStruct with pointer bool",
			generate: func() any {
				return faker.GenerateFake[MinimalStruct]()
			},
			assertions: func(t *testing.T, result any) {
				ms := result.(*MinimalStruct)
				assert.NotNil(t, ms.Flag, "Flag should not be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.generate()
			tt.assertions(t, result)
		})
	}
}

func BenchmarkGenerateFake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = faker.GenerateFake[TestStruct]()
	}
}
