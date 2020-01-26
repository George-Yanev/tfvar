package tfvar

import (
	"bytes"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestLoad(t *testing.T) {
	type args struct {
		rootDir string
	}

	tests := []struct {
		name      string
		args      args
		want      []Variable
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				rootDir: "./testdata/normal",
			},
			want: []Variable{
				{Name: "resource_name"},
				{Name: "instance_name", Value: cty.StringVal("my-instance")},
			},
			assertion: assert.NoError,
		},
		{
			name: "bad",
			args: args{
				rootDir: "./testdata/bad",
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.rootDir)
			tt.assertion(t, err)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestWriteAsEnvVars(t *testing.T) {
	vars, err := Load("testdata/defaults")
	require.NoError(t, err)

	sort.Slice(vars, func(i, j int) bool { return vars[i].Name < vars[j].Name })

	var buf bytes.Buffer
	assert.NoError(t, WriteAsEnvVars(&buf, vars))

	expected := `export TF_VAR_availability_zone_names='["us-west-1a"]'
export TF_VAR_instance_name='my-instance'
`
	assert.Equal(t, expected, buf.String())
}
