package event

import "testing"

func TestRegister(t *testing.T) {
	Init()

	type args struct {
		event   Event
		handler Handler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test-1",
			args: args{
				event: "test-1",
				handler: func(e Event, payload interface{}) {
					if e != "test-1" || payload != "test-1 payload" {
						t.Errorf("received event: %s, data: %v", e, payload)
						return
					}
					t.Logf("received event: %s, data: %v", e, payload)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.args.event, tt.args.handler)
		})
	}

	PublishAsync("test-1", "test-1 payload")
}

func TestUnregister(t *testing.T) {
	Init()

	type args struct {
		event   Event
		handler Handler
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test-1",
			args: args{
				event: "test-1",
				handler: func(e Event, payload interface{}) {
					t.Logf("received event: %s, data: %v", e, payload)
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Unregister(tt.args.event, tt.args.handler)
		})
	}

	PublishSync("test-1", "test-1 payload")
}

func TestPublishAsync(t *testing.T) {
	Init()
	type args struct {
		event   Event
		payload interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test-1",
			args: args{
				event:   "test-1",
				payload: "test-1 payload",
			},
		},
	}

	go func() {
		Register("test-1", func(e Event, payload interface{}) {
			if e != "test-1" || payload != "test-1 payload" {
				t.Errorf("received event: %s, data: %v", e, payload)
				return
			}
			t.Logf("received event: %s, data: %v", e, payload)
		})
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PublishAsync(tt.args.event, tt.args.payload)
		})

	}
}
