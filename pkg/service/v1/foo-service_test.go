package v1

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	v1 "github.com/wingkwong/go-grpc-boilerplate/pkg/api/v1"
)

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
					WithArgs("title", "description", "foo", "foo").
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
	// TO BE IMPLEMENTED
}

func Test_fooServiceServer_Delete(t *testing.T) {
	// TO BE IMPLEMENTED
}
