package repository

import (
	"bytes"
	"context"
	"database/sql"
	"guestbook-example/internal/domain"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initMessageDBMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("failed to open mock sql db, got error: %v", err)
	}

	gormdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		t.Errorf("failed to open gorm db, got error: %v", err)
	}

	return gormdb, mock, db
}

func Test_messageRepo_Create(t *testing.T) {
	buff := &bytes.Buffer{}

	gormdb, mock, db := initMessageDBMock(t)
	defer db.Close()

	type fields struct {
		logger *slog.Logger
		db     *gorm.DB
	}
	type args struct {
		ctx context.Context
		m   *domain.Message
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
				logger: slog.New(slog.NewTextHandler(buff, nil)),
				db: func() *gorm.DB {
					mock.ExpectBegin()
					mock.ExpectExec("INSERT INTO `messages` .*").
						WithArgs(
							sqlmock.AnyArg(),
							sqlmock.AnyArg(),
							nil,
							"Arthur Morgan",
							"Hey, Dutch!",
						).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
					return gormdb
				}(),
			},
			args: args{
				ctx: context.Background(),
				m: &domain.Message{
					Author:  "Arthur Morgan",
					Message: "Hey, Dutch!",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "failed to create message",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				db: func() *gorm.DB {
					mock.ExpectBegin()
					mock.ExpectExec("INSERT INTO `messages` .*").
						WillReturnError(sql.ErrConnDone)
					mock.ExpectRollback()
					return gormdb
				}(),
			},
			args: args{
				ctx: context.Background(),
				m: &domain.Message{
					Author:  "Arthur Morgan",
					Message: "Hey, Dutch!",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MessageRepo{
				logger: tt.fields.logger,
				db:     tt.fields.db,
			}
			got, err := m.Create(tt.args.ctx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("messageRepo.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_messageRepo_Get(t *testing.T) {
	buff := &bytes.Buffer{}

	gormdb, mock, db := initMessageDBMock(t)
	defer db.Close()

	type fields struct {
		logger *slog.Logger
		db     *gorm.DB
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
				logger: slog.New(slog.NewTextHandler(buff, nil)),
				db: func() *gorm.DB {
					mock.ExpectQuery(".*").WithArgs(1).
						WillReturnRows(sqlmock.NewRows([]string{"id", "author", "message"}).
							AddRow(1, "Arthur Morgan", "Hey, Dutch!"))
					return gormdb
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
			wantErr: false,
		},
		{
			name: "failed to get message",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				db: func() *gorm.DB {
					mock.ExpectQuery(".*").WithArgs(1).
						WillReturnError(sql.ErrConnDone)
					return gormdb
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
			m := &MessageRepo{
				logger: tt.fields.logger,
				db:     tt.fields.db,
			}
			got, err := m.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false && got.ID != tt.want.ID {
				t.Errorf("messageRepo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_messageRepo_GetAll(t *testing.T) {
	buff := &bytes.Buffer{}

	gormdb, mock, db := initMessageDBMock(t)
	defer db.Close()

	type fields struct {
		logger *slog.Logger
		db     *gorm.DB
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
				logger: slog.New(slog.NewTextHandler(buff, nil)),
				db: func() *gorm.DB {
					mock.ExpectQuery(".*").
						WillReturnRows(sqlmock.NewRows([]string{"id", "author", "message"}).
							AddRow(1, "Arthur Morgan", "Hey, Dutch!").
							AddRow(2, "Dutch van der Linde", "I have a plan!"))
					return gormdb
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
			},
			wantErr: false,
		},
		{
			name: "get no message",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				db: func() *gorm.DB {
					mock.ExpectQuery(".*").
						WillReturnRows(sqlmock.NewRows([]string{}))
					return gormdb
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []*domain.Message{},
			wantErr: false,
		},
		{
			name: "failed to get all messages",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				db: func() *gorm.DB {
					mock.ExpectQuery(".*").
						WillReturnError(sql.ErrConnDone)
					return gormdb
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
			m := &MessageRepo{
				logger: tt.fields.logger,
				db:     tt.fields.db,
			}
			got, err := m.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepo.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false && len(got) != len(tt.want) {
				t.Errorf("messageRepo.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_messageRepo_Update(t *testing.T) {
	buff := &bytes.Buffer{}

	gormdb, mock, db := initMessageDBMock(t)
	defer db.Close()

	type fields struct {
		logger *slog.Logger
		db     *gorm.DB
	}
	type args struct {
		ctx context.Context
		m   *domain.Message
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
				logger: slog.New(slog.NewTextHandler(buff, nil)),
				db: func() *gorm.DB {
					mock.ExpectBegin()
					mock.ExpectExec("UPDATE `messages` .*").
						WithArgs(
							sqlmock.AnyArg(),
							sqlmock.AnyArg(),
							nil,
							"Arthur Morgan",
							"Hey, Dutch!",
							1,
						).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
					return gormdb
				}(),
			},
			args: args{
				ctx: context.Background(),
				m: &domain.Message{
					ID:      1,
					Author:  "Arthur Morgan",
					Message: "Hey, Dutch!",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "failed to update message",
			fields: fields{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				db: func() *gorm.DB {
					mock.ExpectBegin()
					mock.ExpectExec("UPDATE `messages` .*").
						WillReturnError(sql.ErrConnDone)
					mock.ExpectRollback()
					return gormdb
				}(),
			},
			args: args{
				ctx: context.Background(),
				m: &domain.Message{
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
			m := &MessageRepo{
				logger: tt.fields.logger,
				db:     tt.fields.db,
			}
			if err := m.Update(tt.args.ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("messageRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_messageRepo_Delete(t *testing.T) {
	buff := &bytes.Buffer{}

	gormdb, mock, db := initMessageDBMock(t)
	defer db.Close()

	type fields struct {
		logger *slog.Logger
		db     *gorm.DB
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
				logger: slog.New(slog.NewTextHandler(buff, nil)),
				db: func() *gorm.DB {
					mock.ExpectBegin()
					mock.ExpectExec("UPDATE `messages` .*").
						WithArgs(
							sqlmock.AnyArg(),
							1,
						).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
					return gormdb
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
				db: func() *gorm.DB {
					mock.ExpectBegin()
					mock.ExpectExec("UPDATE `messages` .*").
						WillReturnError(sql.ErrConnDone)
					mock.ExpectRollback()
					return gormdb
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
			m := &MessageRepo{
				logger: tt.fields.logger,
				db:     tt.fields.db,
			}
			if err := m.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("messageRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
