package utils

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func BindProtoJSON(c *gin.Context, msg proto.Message) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	if len(body) == 0 {
		return fmt.Errorf("request body is empty")
	}

	opts := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}
	if err := opts.Unmarshal(body, msg); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	return nil
}
