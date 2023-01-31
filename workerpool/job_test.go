package workerpool_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/fgmaia/goworkerpool/workerpool"
)

var (
	errDefault = errors.New("wrong argument type")

	descriptor = workerpool.JobDescriptor{
		ID:    workerpool.JobID("1"),
		JType: workerpool.JobType("anyType"),
		Metadata: workerpool.JobMetadata{
			"foo": "foo",
			"bar": "bar",
		},
	}

	execFn = func(ctx context.Context, args interface{}) (interface{}, error) {
		argVal, ok := args.(int)
		if !ok {
			return nil, errDefault
		}

		return argVal * 2, nil
	}
)

func Test_job_Execute(t *testing.T) {
	ctx := context.TODO()

	type fields struct {
		descriptor workerpool.JobDescriptor
		execFn     workerpool.ExecutionFn
		args       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   workerpool.Result
	}{
		{
			name: "job execution success",
			fields: fields{
				descriptor: descriptor,
				execFn:     execFn,
				args:       10,
			},
			want: workerpool.Result{
				Value:      20,
				Descriptor: descriptor,
			},
		},
		{
			name: "job execution failure",
			fields: fields{
				descriptor: descriptor,
				execFn:     execFn,
				args:       "10",
			},
			want: workerpool.Result{
				Err:        errDefault,
				Descriptor: descriptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := workerpool.NewJob(
				tt.fields.descriptor,
				tt.fields.execFn,
				tt.fields.args,
			)

			gotResult := j.Execute(ctx)
			if tt.want.Err != nil {
				if !reflect.DeepEqual(gotResult.Err.Error(), tt.want.Err.Error()) {
					t.Errorf("execute() = %v, wantError %v", gotResult.Err, tt.want.Err)
				}
				return
			}

			if !reflect.DeepEqual(gotResult, tt.want) {
				t.Errorf("execute() = %v, want %v", gotResult, tt.want)
			}
		})
	}
}
