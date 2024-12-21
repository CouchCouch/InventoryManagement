package middleware

import (
    "errors"
    "net/http"

    "inventoryapi/api"
    "inventoryapi/internal/tools"
    log "github.com/sirupsen/logrus"
)

var UnAuthorizedError = errors.New("Invalid Username or Token.")

func Authorization(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        var user string = r.URL.Query().Get("username")
        var token = r.Header.Get("authorization")
        var err error

        if user == "" {
        }

        if user == "" || token == "" {
            log.Error(UnAuthorizedError)
            api.RequestErrorHandler(w, UnAuthorizedError)
            return
        }

        var database *tools.DatabaseInterface
        database, err = tools.NewDatabase()
        if err != nil {
            api.InternalErrorHandler(w)
            return
        }

        var loginDetails *tools.LoginDetails
        loginDetails = (*database).GetUserLoginDetails(user)

        if(loginDetails == nil || (token != (*loginDetails).AuthToken)) {
            log.Error(UnAuthorizedError)
            api.RequestErrorHandler(w, UnAuthorizedError)
            return
        }

        next.ServeHTTP(w, r)
    })
}
