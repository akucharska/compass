package tenantmapping

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kyma-incubator/compass/components/director/internal/persistence"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//go:generate mockery -name=ScopesGetter -output=automock -outpkg=automock -case=underscore
type ScopesGetter interface {
	GetRequiredScopes(scopesDefinition string) ([]string, error)
}

//go:generate mockery -name=ReqDataParser -output=automock -outpkg=automock -case=underscore
type ReqDataParser interface {
	Parse(req *http.Request) (ReqData, error)
}

//go:generate mockery -name=ObjectContextForUserProvider -output=automock -outpkg=automock -case=underscore
type ObjectContextForUserProvider interface {
	GetObjectContext(reqData ReqData, authID string) (ObjectContext, error)
}

//go:generate mockery -name=ObjectContextForSystemAuthProvider -output=automock -outpkg=automock -case=underscore
type ObjectContextForSystemAuthProvider interface {
	GetObjectContext(ctx context.Context, reqData ReqData, authID string, authFlow AuthFlow) (ObjectContext, error)
}

type Handler struct {
	reqDataParser       ReqDataParser
	transact            persistence.Transactioner
	mapperForUser       ObjectContextForUserProvider
	mapperForSystemAuth ObjectContextForSystemAuthProvider
}

func NewHandler(
	reqDataParser ReqDataParser,
	transact persistence.Transactioner,
	mapperForUser ObjectContextForUserProvider,
	mapperForSystemAuth ObjectContextForSystemAuthProvider) *Handler {
	return &Handler{
		reqDataParser:       reqDataParser,
		transact:            transact,
		mapperForUser:       mapperForUser,
		mapperForSystemAuth: mapperForSystemAuth,
	}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(writer, fmt.Sprintf("Bad request method. Got %s, expected POST", req.Method), http.StatusBadRequest)
		return
	}

	reqData, err := h.reqDataParser.Parse(req)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err, "while parsing the request")
		return
	}

	tx, err := h.transact.Begin()
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err, "while opening the db transaction")
		return
	}
	defer h.transact.RollbackUnlessCommited(tx)

	ctx := persistence.SaveToContext(req.Context(), tx)

	objCtx, err := h.getObjectContext(ctx, reqData)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err, "while getting object context")
		return
	}

	reqData.Body.Extra["tenant"] = objCtx.TenantID
	reqData.Body.Extra["scope"] = objCtx.Scopes
	reqData.Body.Extra["objectID"] = objCtx.ObjectID
	reqData.Body.Extra["objectType"] = objCtx.ObjectType

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reqData.Body)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err, "while encoding data")
		return
	}
}

func (h *Handler) getObjectContext(ctx context.Context, reqData ReqData) (ObjectContext, error) {
	authID, authFlow, err := reqData.GetAuthID()
	if err != nil {
		return ObjectContext{}, errors.Wrap(err, "while determining the auth ID from the request")
	}

	switch authFlow {
	case JWTAuthFlow:
		return h.mapperForUser.GetObjectContext(reqData, authID)
	case OAuth2Flow, CertificateFlow:
		return h.mapperForSystemAuth.GetObjectContext(ctx, reqData, authID, authFlow)
	}

	return ObjectContext{}, fmt.Errorf("unknown authentication flow (%s)", authFlow)
}

func respondWithError(writer http.ResponseWriter, httpErrorCode int, err error, wrapperStr string) {
	wrappedErr := errors.Wrap(err, wrapperStr)
	log.Error(wrappedErr)

	http.Error(writer, wrappedErr.Error(), httpErrorCode)
}
