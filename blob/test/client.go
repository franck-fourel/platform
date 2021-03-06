package test

import (
	"context"

	"github.com/tidepool-org/platform/blob"
	"github.com/tidepool-org/platform/page"
)

type ListInput struct {
	Context    context.Context
	UserID     string
	Filter     *blob.Filter
	Pagination *page.Pagination
}

type ListOutput struct {
	Blobs blob.Blobs
	Error error
}

type CreateInput struct {
	Context context.Context
	UserID  string
	Create  *blob.Create
}

type CreateOutput struct {
	Blob  *blob.Blob
	Error error
}

type GetInput struct {
	Context context.Context
	ID      string
}

type GetOutput struct {
	Blob  *blob.Blob
	Error error
}

type GetContentInput struct {
	Context context.Context
	ID      string
}

type GetContentOutput struct {
	Content *blob.Content
	Error   error
}

type DeleteInput struct {
	Context context.Context
	ID      string
}

type DeleteOutput struct {
	Deleted bool
	Error   error
}

type Client struct {
	ListInvocations       int
	ListInputs            []ListInput
	ListStub              func(ctx context.Context, userID string, filter *blob.Filter, pagination *page.Pagination) (blob.Blobs, error)
	ListOutputs           []ListOutput
	ListOutput            *ListOutput
	CreateInvocations     int
	CreateInputs          []CreateInput
	CreateStub            func(ctx context.Context, userID string, create *blob.Create) (*blob.Blob, error)
	CreateOutputs         []CreateOutput
	CreateOutput          *CreateOutput
	GetInvocations        int
	GetInputs             []GetInput
	GetStub               func(ctx context.Context, id string) (*blob.Blob, error)
	GetOutputs            []GetOutput
	GetOutput             *GetOutput
	GetContentInvocations int
	GetContentInputs      []GetContentInput
	GetContentStub        func(ctx context.Context, id string) (*blob.Content, error)
	GetContentOutputs     []GetContentOutput
	GetContentOutput      *GetContentOutput
	DeleteInvocations     int
	DeleteInputs          []DeleteInput
	DeleteStub            func(ctx context.Context, id string) (bool, error)
	DeleteOutputs         []DeleteOutput
	DeleteOutput          *DeleteOutput
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) List(ctx context.Context, userID string, filter *blob.Filter, pagination *page.Pagination) (blob.Blobs, error) {
	c.ListInvocations++
	c.ListInputs = append(c.ListInputs, ListInput{Context: ctx, UserID: userID, Filter: filter, Pagination: pagination})
	if c.ListStub != nil {
		return c.ListStub(ctx, userID, filter, pagination)
	}
	if len(c.ListOutputs) > 0 {
		output := c.ListOutputs[0]
		c.ListOutputs = c.ListOutputs[1:]
		return output.Blobs, output.Error
	}
	if c.ListOutput != nil {
		return c.ListOutput.Blobs, c.ListOutput.Error
	}
	panic("List has no output")
}

func (c *Client) Create(ctx context.Context, userID string, create *blob.Create) (*blob.Blob, error) {
	c.CreateInvocations++
	c.CreateInputs = append(c.CreateInputs, CreateInput{Context: ctx, UserID: userID, Create: create})
	if c.CreateStub != nil {
		return c.CreateStub(ctx, userID, create)
	}
	if len(c.CreateOutputs) > 0 {
		output := c.CreateOutputs[0]
		c.CreateOutputs = c.CreateOutputs[1:]
		return output.Blob, output.Error
	}
	if c.CreateOutput != nil {
		return c.CreateOutput.Blob, c.CreateOutput.Error
	}
	panic("Create has no output")
}

func (c *Client) Get(ctx context.Context, id string) (*blob.Blob, error) {
	c.GetInvocations++
	c.GetInputs = append(c.GetInputs, GetInput{Context: ctx, ID: id})
	if c.GetStub != nil {
		return c.GetStub(ctx, id)
	}
	if len(c.GetOutputs) > 0 {
		output := c.GetOutputs[0]
		c.GetOutputs = c.GetOutputs[1:]
		return output.Blob, output.Error
	}
	if c.GetOutput != nil {
		return c.GetOutput.Blob, c.GetOutput.Error
	}
	panic("Get has no output")
}

func (c *Client) GetContent(ctx context.Context, id string) (*blob.Content, error) {
	c.GetContentInvocations++
	c.GetContentInputs = append(c.GetContentInputs, GetContentInput{Context: ctx, ID: id})
	if c.GetContentStub != nil {
		return c.GetContentStub(ctx, id)
	}
	if len(c.GetContentOutputs) > 0 {
		output := c.GetContentOutputs[0]
		c.GetContentOutputs = c.GetContentOutputs[1:]
		return output.Content, output.Error
	}
	if c.GetContentOutput != nil {
		return c.GetContentOutput.Content, c.GetContentOutput.Error
	}
	panic("GetContent has no output")
}

func (c *Client) Delete(ctx context.Context, id string) (bool, error) {
	c.DeleteInvocations++
	c.DeleteInputs = append(c.DeleteInputs, DeleteInput{Context: ctx, ID: id})
	if c.DeleteStub != nil {
		return c.DeleteStub(ctx, id)
	}
	if len(c.DeleteOutputs) > 0 {
		output := c.DeleteOutputs[0]
		c.DeleteOutputs = c.DeleteOutputs[1:]
		return output.Deleted, output.Error
	}
	if c.DeleteOutput != nil {
		return c.DeleteOutput.Deleted, c.DeleteOutput.Error
	}
	panic("Delete has no output")
}

func (c *Client) AssertOutputsEmpty() {
	if len(c.ListOutputs) > 0 {
		panic("ListOutputs is not empty")
	}
	if len(c.CreateOutputs) > 0 {
		panic("CreateOutputs is not empty")
	}
	if len(c.GetOutputs) > 0 {
		panic("GetOutputs is not empty")
	}
	if len(c.GetContentOutputs) > 0 {
		panic("GetContentOutputs is not empty")
	}
	if len(c.DeleteOutputs) > 0 {
		panic("DeleteOutputs is not empty")
	}
}
