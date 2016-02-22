package user

import (
	"net/http"

	log "github.com/tidepool-org/platform/logger"

	"github.com/tidepool-org/platform/Godeps/_workspace/src/github.com/ant0ine/go-json-rest/rest"
)

type ChainedMiddleware func(rest.HandlerFunc) rest.HandlerFunc

//Authorization middleware is used for validation of incoming tokens
type AuthorizationMiddleware struct {
	Client Client
}

func NewAuthorizationMiddleware(userClient Client) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{Client: userClient}
}

//Valid - then we continue
//Invalid - then we return 401 (http.StatusUnauthorized)
func (mw *AuthorizationMiddleware) ValidateToken(h rest.HandlerFunc) rest.HandlerFunc {

	return func(w rest.ResponseWriter, r *rest.Request) {

		token := r.Header.Get(x_tidepool_session_token)
		userid := r.PathParam("userid")

		if tokenData := mw.Client.CheckToken(token); tokenData != nil {
			if tokenData.IsServer || tokenData.UserID == userid {
				h(w, r)
				return
			}
			log.Logging.Info("id's don't match and not server token", tokenData.UserID, userid)
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

//Authorization middleware is used for getting user permissons
type PermissonsMiddleware struct {
	Client Client
}

const PERMISSIONS = "PERMISSIONS"

func NewPermissonsMiddleware(userClient Client) *PermissonsMiddleware {
	return &PermissonsMiddleware{Client: userClient}
}

//Attach permissons if they exist
//http.StatusInternalServerError if there is an error getting the user permissons
func (mw *PermissonsMiddleware) GetPermissons(h rest.HandlerFunc) rest.HandlerFunc {

	return func(w rest.ResponseWriter, r *rest.Request) {

		token := r.Header.Get(x_tidepool_session_token)
		userid := r.PathParam("userid")

		permissions, err := mw.Client.GetUserPermissons(userid, token)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Env[PERMISSIONS] = permissions
		h(w, r)
		return

	}
}
