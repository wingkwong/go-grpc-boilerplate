package v1

import (
	"context"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	v1 "github.com/wingkwong/go-grpc-boilerplate/pkg/api/v1"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func Test_fooServiceServer_Create(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("[Error] '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewFooServiceServer(db)

	type args struct {
		ctx context.Context
		req *v1.CreateRequest
	}
	tests := []struct {
		name    string
		s       v1.FooServiceServer
		args    args
		mock    func()
		want    *v1.CreateResponse
		wantErr bool
	}{
		{
			name: "01 - OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					ApiVersion: "v1",
					Foo: &v1.Foo{
						Title: "title",
						Desc:  "description",
						SysFields: &v1.SystemFields{
							CreatedBy: "foo",
							UpdatedBy: "foo",
						},
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO Foo").
					WithArgs("title", "description", "foo", "foo", AnyTime{}, AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &v1.CreateResponse{
				ApiVersion: "v1",
				Id:         1,
			},
		},
		{
			name: "02 - Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					ApiVersion: "v1000",
					Foo: &v1.Foo{
						Title: "title",
						Desc:  "description",
						SysFields: &v1.SystemFields{
							CreatedBy: "foo",
							UpdatedBy: "foo",
						},
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "03 - INSERT failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					ApiVersion: "v1",
					Foo: &v1.Foo{
						Title: "title",
						Desc:  "description",
						SysFields: &v1.SystemFields{
							CreatedBy: "foo",
							UpdatedBy: "foo",
						},
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO Foo").
					WithArgs("title", "description", "foo", "foo").
					WillReturnError(errors.New("INSERT failed"))
			},
			wantErr: true,
		},
		{
			name: "04 - LastInsertId failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					ApiVersion: "v1",
					Foo: &v1.Foo{
						Title: "title",
						Desc:  "description",
						SysFields: &v1.SystemFields{
							CreatedBy: "foo",
							UpdatedBy: "foo",
						},
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO Foo").
					WithArgs("title", "description", "foo", "foo").
					WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId failed")))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("fooServiceServer.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fooServiceServer.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fooServiceServer_Read(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("[Error] '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewFooServiceServer(db)

	curTime := time.Now()

	type args struct {
		ctx context.Context
		req *v1.ReadRequest
	}
	tests := []struct {
		name    string
		s       v1.FooServiceServer
		args    args
		mock    func()
		want    *v1.ReadResponse
		wantErr bool
	}{
		{
			name: "01 - OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					ApiVersion: "v1",
					Id:         1,
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"ID", "Title", "Desc", "CreatedBy", "UpdatedBy", "CreatedAt", "UpdatedAt"}).
					AddRow(1, "title", "description", "foo", "foo", curTime, curTime)
				mock.ExpectQuery("SELECT (.+) FROM Foo").WithArgs(1).WillReturnRows(rows)
			},
			want: &v1.ReadResponse{
				ApiVersion: "v1",
				Foo: &v1.Foo{
					Id:    1,
					Title: "title",
					Desc:  "description",
					SysFields: &v1.SystemFields{
						CreatedBy: "foo",
						UpdatedBy: "foo",
						CreatedAt: timestamppb.New(curTime),
						UpdatedAt: timestamppb.New(curTime),
					},
				},
			},
		},
		{
			name: "02 - Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					ApiVersion: "v1000",
					Id:         1,
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "03 - SELECT failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					ApiVersion: "v1",
					Id:         1,
				},
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM foo").WithArgs(1).
					WillReturnError(errors.New("SELECT Failed"))
			},
			wantErr: true,
		},
		{
			name: "04 - Not found",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					ApiVersion: "v1",
					Id:         1,
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"ID", "Title", "Desc", "CreatedBy", "UpdatedBy", "CreatedAt", "UpdatedAt"})
				mock.ExpectQuery("SELECT (.+) FROM Foo").WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Read(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("fooServiceServer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fooServiceServer.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fooServiceServer_Update(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("[Error] '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewFooServiceServer(db)

	type args struct {
		ctx context.Context
		req *v1.UpdateRequest
	}
	tests := []struct {
		name    string
		s       v1.FooServiceServer
		args    args
		mock    func()
		want    *v1.UpdateResponse
		wantErr bool
	}{
		{
			name: "01 - OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					ApiVersion: "v1",
					Foo: &v1.Foo{
						Id:    1,
						Title: "new title",
						Desc:  "new description",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE Foo").WithArgs("new title", "new description", AnyTime{}, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &v1.UpdateResponse{
				ApiVersion: "v1",
				Count:      1,
			},
		},
		{
			name: "02 - Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					ApiVersion: "v1",
					Foo: &v1.Foo{
						Id:    1,
						Title: "new title",
						Desc:  "new description",
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "03 - UPDATE failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					ApiVersion: "v1",
					Foo: &v1.Foo{
						Id:    1,
						Title: "new title",
						Desc:  "new description",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE Foo").WithArgs("new title", "new description", 1).
					WillReturnError(errors.New("UPDATE failed"))
			},
			wantErr: true,
		},
		{
			name: "04 - RowsAffected failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					ApiVersion: "v1",
					Foo: &v1.Foo{
						Id:    1,
						Title: "new title",
						Desc:  "new description",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE Foo").WithArgs("new title", "new description", 1).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("RowsAffected failed")))
			},
			wantErr: true,
		},
		{
			name: "05 - Not Found",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					ApiVersion: "v1",
					Foo: &v1.Foo{
						Id:    1,
						Title: "new title",
						Desc:  "new description",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE Foo").WithArgs("new title", "new description", 1).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("fooServiceServer.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fooServiceServer.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fooServiceServer_Delete(t *testing.T) {
	// TO BE IMPLEMENTED
}
