package user

import (
	"github.com/reearth/reearthx/account/accountdomain"
)

type ID = accountdomain.UserID
type IDList = accountdomain.UserIDList
type WorkspaceID = accountdomain.WorkspaceID
type WorkspaceIDList = accountdomain.WorkspaceIDList

var NewID = accountdomain.NewUserID
var NewWorkspaceID = accountdomain.NewWorkspaceID

var IDFrom = accountdomain.UserIDFrom
var WorkspaceIDFrom = accountdomain.WorkspaceIDFrom

var IDFromRef = accountdomain.UserIDFromRef
var WorkspaceIDFromRef = accountdomain.WorkspaceIDFromRef

var ErrInvalidID = accountdomain.ErrInvalidID
