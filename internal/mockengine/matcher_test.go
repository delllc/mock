package mockengine

import "testing"

func TestMatchPath(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		{
			name:    "exact match",
			pattern: "/api/products",
			path:    "/api/products",
			want:    true,
		},
		{
			name:    "parameter match",
			pattern: "/api/products/:id",
			path:    "/api/products/42",
			want:    true,
		},
		{
			name:    "different parameter values",
			pattern: "/api/products/:id",
			path:    "/api/products/999",
			want:    true,
		},
		{
			name:    "extra path segment",
			pattern: "/api/products/:id",
			path:    "/api/products/42/reviews",
			want:    false,
		},
		{
			name:    "missing parameter",
			pattern: "/api/products/:id",
			path:    "/api/products",
			want:    false,
		},
		{
			name:    "empty parameter",
			pattern: "/api/products/:id",
			path:    "/api/products/",
			want:    false,
		},
		{
			name:    "different static segment",
			pattern: "/api/products/:id",
			path:    "/api/users/42",
			want:    false,
		},
		{
			name:    "root path",
			pattern: "/",
			path:    "/",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchPath(tt.pattern, tt.path)

			if got != tt.want {
				t.Errorf(
					"matchPath(%q, %q) = %v, want %v",
					tt.pattern,
					tt.path,
					got,
					tt.want,
				)
			}
		})
	}
}
