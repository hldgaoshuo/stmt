package compiler

import "testing"

func Test_objectsEmit(t *testing.T) {
	tests := []struct {
		name    string
		objects []*Object
	}{
		{
			name: "1",
			objects: []*Object{
				{
					Literal:    int64(1),
					ObjectType: INT,
				},
				{
					Literal:    int64(2),
					ObjectType: INT,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objectsEmit(tt.objects)
		})
	}
}
