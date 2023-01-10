package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func ResourceParseId(ctx context.Context, id string) (string, string, error) {
	tflog.Debug(ctx, fmt.Sprintf("ID to be parsed: %s", id))
	unquotedId, err := strconv.Unquote(id)
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("cannot unquote id, skip unquoting: %s: %s", id, err))
		unquotedId = id
	}
	tflog.Debug(ctx, fmt.Sprintf("ID to be parsed: %s", id))
  parts := strings.SplitN(unquotedId, "/", 2)

  if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
    return "", "", fmt.Errorf("unexpected format of ID (%s), expected attribute1/attribute2", id)
  }
	tflog.Debug(ctx, fmt.Sprintf("Parsed ID is: %v", parts))
  return parts[0], parts[1], nil
}

func CreateResourceId(parent string, leaf string) string {
	return strings.Join([]string{parent, leaf}, "/")
}

// RESOURCES CHANGE PLAN for I64

func Int64DefaultValue(v types.Int64) planmodifier.Int64 {
	return &Int64DefaultModifier{v}
}

type Int64DefaultModifier struct {
	DefaultValue types.Int64
}

var _ planmodifier.Int64 = (*Int64DefaultModifier)(nil)

func (apm *Int64DefaultModifier) Description(ctx context.Context) string {
	return ""
}
func (apm *Int64DefaultModifier) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (m *Int64DefaultModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// If the value is unknown or known, do not set default value.
	if !req.PlanValue.IsNull() {
			return
	}
	resp.PlanValue = m.DefaultValue
}

// RESOURCES CHANGE PLAN for STRING

func StringDefaultValue(v types.String) planmodifier.String {
	return &StringDefaultModifier{v}
}

type StringDefaultModifier struct {
	DefaultValue types.String
}

var _ planmodifier.String = (*StringDefaultModifier)(nil)

func (apm *StringDefaultModifier) Description(ctx context.Context) string {
	return ""
}
func (apm *StringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (m *StringDefaultModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the value is unknown or known, do not set default value.
	if !req.PlanValue.IsNull() {
			return
	}
	resp.PlanValue = m.DefaultValue
}
