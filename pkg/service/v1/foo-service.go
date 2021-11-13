package v1

import (
	"context"
	"database/sql"
	"time"

	v1 "github.com/wingkwong/go-grpc-boilerplate/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	apiVersion = "v1"
)

type fooServiceServer struct {
	db *sql.DB
}

func NewFooServiceServer(db *sql.DB) v1.FooServiceServer {
	return &fooServiceServer{db: db}
}

func (s *fooServiceServer) checkAPI(api string) error {
	if len(api) > 0 && apiVersion != api {
		return status.Errorf(codes.Unimplemented,
			"[Error] Unsupported API version: service API version '%s', but got '%s'", apiVersion, api)
	}
	return nil
}

func (s *fooServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "[Error] Failed to connect to database: "+err.Error())
	}
	return c, nil
}

func (s *fooServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := s.checkAPI(req.ApiVersion); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(
		ctx,
		"INSERT INTO Foo(`Title`, `Desc`, `CreatedBy`, `UpdatedBy`) VALUES(?, ?, ?, ?)",
		req.Foo.Title, req.Foo.Desc, req.Foo.SysFields.CreatedBy, req.Foo.SysFields.UpdatedBy)
	if err != nil {
		return nil, status.Error(codes.Unknown, "[Error] Failed to insert into record: "+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "[Error] Failed to retrieve last inserted id:  "+err.Error())
	}

	return &v1.CreateResponse{
		ApiVersion: apiVersion,
		Id:         id,
	}, nil
}

func (s *fooServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	if err := s.checkAPI(req.ApiVersion); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	id := req.Id
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Title`, `Desc`, `CreatedBy`, `UpdatedBy`, `CreatedAt`, `UpdatedAt` FROM Foo WHERE `ID` = ?", id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "[Error] Failed to select data from Foo by Id %d : "+err.Error(), id)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "[Error] Failed to retrieve data from Foo: "+err.Error())
		}
		return nil, status.Errorf(codes.NotFound, "[Error] Failed to find Id : %d", id)
	}

	var foo v1.Foo
	foo.SysFields = &v1.SystemFields{}

	var CreatedAt time.Time
	var UpdatedAt time.Time

	// TODO: probably use jmoiron/sqlx to assign to a struct
	if err := rows.Scan(&foo.Id, &foo.Title, &foo.Desc, &foo.SysFields.CreatedBy, &foo.SysFields.UpdatedBy, &CreatedAt, &UpdatedAt); err != nil {
		return nil, status.Error(codes.Unknown, "[Error] Failed to retrieve values from Foo rows : "+err.Error())
	}

	foo.SysFields.CreatedAt = timestamppb.New(CreatedAt)
	foo.SysFields.UpdatedAt = timestamppb.New(UpdatedAt)

	if rows.Next() {
		return nil, status.Errorf(codes.Unknown, "[Error] multiple rows with the same id :'%d'", id)
	}

	return &v1.ReadResponse{
		ApiVersion: apiVersion,
		Foo:        &foo,
	}, nil
}

func (s *fooServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := s.checkAPI(req.ApiVersion); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Title`, `Desc`, `CreatedBy`, `UpdatedBy`, `CreatedAt`, `UpdatedAt` FROM Foo")
	if err != nil {
		return nil, status.Error(codes.Unknown, "[Error] Failed to retrieve all data from Foo: "+err.Error())
	}
	defer rows.Close()

	var CreatedAt time.Time
	var UpdatedAt time.Time

	fooList := []*v1.Foo{}
	for rows.Next() {
		foo := new(v1.Foo)
		if err := rows.Scan(&foo.Id, &foo.Title, &foo.Desc, &foo.SysFields.CreatedBy, &foo.SysFields.UpdatedBy, &CreatedAt, &UpdatedAt); err != nil {
			return nil, status.Error(codes.Unknown, "[ERROR] Failed to retrieve field values from Foo: "+err.Error())
		}

		foo.SysFields.CreatedAt = timestamppb.New(CreatedAt)
		if err != nil {
			return nil, status.Error(codes.Unknown, "[ERROR] Field createdAt has invalid format: "+err.Error())
		}

		foo.SysFields.UpdatedAt = timestamppb.New(UpdatedAt)
		if err != nil {
			return nil, status.Error(codes.Unknown, "[ERROR] Field updatedAt has invalid format: "+err.Error())
		}
		fooList = append(fooList, foo)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "[ERROR] Failed to retrieve data from Foo: "+err.Error())
	}

	return &v1.ReadAllResponse{
		ApiVersion: apiVersion,
		Foos:       fooList,
	}, nil
}

func (s *fooServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	if err := s.checkAPI(req.ApiVersion); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "UPDATE Foo SET `Title` = ?, `Desc` = ? WHERE `ID` = ?", req.Foo.Title, req.Foo.Desc, req.Foo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "[Error] Failed to update Foo : "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "[Error] Failed to retrieve rows affected value :  "+err.Error())
	}

	if rows == 0 {
		return nil, status.Errorf(codes.NotFound, "[Error] Failed to update Foo with id : %d", req.Foo.Id)
	}

	return &v1.UpdateResponse{
		ApiVersion: apiVersion,
		Count:      rows,
	}, nil
}

func (s *fooServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := s.checkAPI(req.ApiVersion); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	id := req.Id

	res, err := c.ExecContext(ctx, "DELETE FROM Foo WHERE `ID` = ?", id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "[Error] Failed to delete Foo : "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "[Error] Failed to retrieve rows affected value :  "+err.Error())
	}

	if rows == 0 {
		return nil, status.Errorf(codes.NotFound, "[Error] Failed to delete Foo with id : %d", id)
	}

	return &v1.DeleteResponse{
		ApiVersion: apiVersion,
		Count:      rows,
	}, nil
}
