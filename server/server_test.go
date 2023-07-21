package server

import (
	"github.com/RaphaelL2e/Rasync/broker"
	"testing"
)

func TestServer_Send(t *testing.T) {
	broker1, err := broker.NewBroker("localhost", "6379", "")
	if err != nil {
		panic(err)
	}
	type fields struct {
		groupName string
		broker    *broker.Broker
	}
	type args struct {
		worker string
		data   interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantJobId string
		wantErr   bool
	}{
		{
			name: "test",
			fields: fields{
				groupName: "test",
				broker:    broker1,
			},
			args: args{
				worker: "worker1",
				data:   "hello",
			},
			wantErr: false,
		},
		{
			name: "test",
			fields: fields{
				groupName: "test",
				broker:    broker1,
			},
			args: args{
				worker: "worker1",
				data:   "world",
			},
			wantErr: false,
		},
		{
			name: "test",
			fields: fields{
				groupName: "test",
				broker:    broker1,
			},
			args: args{
				worker: "worker1",
				data:   "!!!",
			},
			wantErr: false,
		},
		{
			name: "test",
			fields: fields{
				groupName: "test",
				broker:    broker1,
			},
			args: args{
				worker: "worker1",
				data:   "123",
			},
			wantErr: false,
		},
		{
			name: "test",
			fields: fields{
				groupName: "test",
				broker:    broker1,
			},
			args: args{
				worker: "worker1",
				data:   "456",
			},
			wantErr: false,
		},
		{
			name: "test",
			fields: fields{
				groupName: "test",
				broker:    broker1,
			},
			args: args{
				worker: "worker1",
				data:   "789",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				groupName: tt.fields.groupName,
				broker:    tt.fields.broker,
			}
			jobId, err := s.Send(tt.args.worker, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(jobId)
		})
	}
}

func TestServer_Cancel(t *testing.T) {
	broker1, err := broker.NewBroker("localhost", "6379", "")
	if err != nil {
		panic(err)
	}
	type fields struct {
		groupName string
		broker    *broker.Broker
	}
	type args struct {
		worker string
		jobId  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				groupName: "test",
				broker:    broker1,
			},
			args: args{
				worker: "worker1",
				jobId:  "e62bc823-8787-4ef5-9aa8-795b45d345dc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				groupName: tt.fields.groupName,
				broker:    tt.fields.broker,
			}
			if err := s.Cancel(tt.args.worker, tt.args.jobId); (err != nil) != tt.wantErr {
				t.Errorf("Cancel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
