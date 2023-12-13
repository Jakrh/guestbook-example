package service

import (
	"context"
	"fmt"
	"guestbook-example/internal/domain"
	"guestbook-example/internal/service/mocks"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestMessageService_Get(t *testing.T) {
	type fields struct {
		logger      *slog.Logger
		messageRepo MessageRepo
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Message
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					mockRepo.On("Get", mock.Anything, mock.AnythingOfType("int64")).Return(
						&domain.Message{
							ID:      1,
							Author:  "Arthur Morgan",
							Message: "Hey, Dutch!",
						}, nil)
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &domain.Message{
				ID:      1,
				Author:  "Arthur Morgan",
				Message: "Hey, Dutch!",
			},
		},
		{
			name: "failed to get message",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					mockRepo.On("Get", mock.Anything, mock.AnythingOfType("int64")).Return(
						nil,
						fmt.Errorf("failed to get message"))
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMessageService(tt.fields.logger, tt.fields.messageRepo)
			got, err := s.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MessageService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageService_GetAll(t *testing.T) {
	type fields struct {
		logger      *slog.Logger
		messageRepo MessageRepo
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*domain.Message
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					mockRepo.On("GetAll", mock.Anything).Return(
						[]*domain.Message{
							{
								ID:      1,
								Author:  "Arthur Morgan",
								Message: "Hey, Dutch!",
							},
							{
								ID:      2,
								Author:  "Dutch van der Linde",
								Message: "I have a plan!",
							},
							{
								ID:      3,
								Author:  "John Marston",
								Message: "I have a family!",
							},
						}, nil)
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			want: []*domain.Message{
				{
					ID:      1,
					Author:  "Arthur Morgan",
					Message: "Hey, Dutch!",
				},
				{
					ID:      2,
					Author:  "Dutch van der Linde",
					Message: "I have a plan!",
				},
				{
					ID:      3,
					Author:  "John Marston",
					Message: "I have a family!",
				},
			},
		},
		{
			name: "failed to get all messages",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					mockRepo.On("GetAll", mock.Anything).Return(
						nil,
						fmt.Errorf("failed to get all messages"))
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMessageService(tt.fields.logger, tt.fields.messageRepo)
			got, err := s.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MessageService.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageService_Create(t *testing.T) {
	type fields struct {
		logger      *slog.Logger
		messageRepo MessageRepo
	}
	type args struct {
		ctx     context.Context
		message *domain.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(
						int64(1), nil)
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
				message: &domain.Message{
					ID:      1,
					Author:  "Arthur Morgan",
					Message: "Hey, Dutch!",
				},
			},
			want: int64(1),
		},
		{
			name: "failed to create message",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(
						int64(0), fmt.Errorf("failed to create message"))
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
				message: &domain.Message{
					ID:      1,
					Author:  "Arthur Morgan",
					Message: "Hey, Dutch!",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMessageService(tt.fields.logger, tt.fields.messageRepo)
			got, err := s.Create(tt.args.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MessageService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageService_Update(t *testing.T) {
	type fields struct {
		logger      *slog.Logger
		messageRepo MessageRepo
	}
	type args struct {
		ctx     context.Context
		message *domain.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(
						nil)
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
				message: &domain.Message{
					ID:      1,
					Author:  "Arthur Morgan",
					Message: "Hey, Dutch!",
				},
			},
			wantErr: false,
		},
		{
			name: "failed to update message",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(
						fmt.Errorf("failed to update message"))
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
				message: &domain.Message{
					ID:      1,
					Author:  "Arthur Morgan",
					Message: "Hey, Dutch!",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMessageService(tt.fields.logger, tt.fields.messageRepo)
			if err := s.Update(tt.args.ctx, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("MessageService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageService_Delete(t *testing.T) {
	type fields struct {
		logger      *slog.Logger
		messageRepo MessageRepo
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					mockRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(
						nil)
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: false,
		},
		{
			name: "failed to delete message",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					mockRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(
						fmt.Errorf("failed to delete message"))
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMessageService(tt.fields.logger, tt.fields.messageRepo)
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("MessageService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMessageService(t *testing.T) {
	type args struct {
		logger      *slog.Logger
		messageRepo MessageRepo
	}
	tests := []struct {
		name string
		args args
		want *MessageService
	}{
		{
			name: "success",
			args: args{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					return mockRepo
				}(),
			},
			want: &MessageService{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				messageRepo: func() MessageRepo {
					mockRepo := new(mocks.MessageRepo)
					return mockRepo
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			if got := NewMessageService(logger, tt.args.messageRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessageService() = %v, want %v", got, tt.want)
			}
		})
	}
}
