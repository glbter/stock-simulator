package lambda

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/glbter/currency-ex/currency/exchanger/factory"
	"github.com/glbter/currency-ex/pkg/serrors"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
	pgx2 "github.com/glbter/currency-ex/pkg/sql/pgx"
	"github.com/glbter/currency-ex/stocks/repository/postgres"
	"github.com/glbter/currency-ex/stocks/usecases"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
)

const (
	DsnEnv = "DSN"

	DBSchema = "public"
)

type UserIDExtractor interface {
	GetUserID(request events.APIGatewayV2HTTPRequest) (string, error)
}

var _ UserIDExtractor = &UserIDExtractorFromAuthHeader{}

type UserIDExtractorFromAuthHeader struct {
	parser JwtParser
}

func (h UserIDExtractorFromAuthHeader) GetUserID(request events.APIGatewayV2HTTPRequest) (string, error) {
	t, ok := request.Headers["Authorization"]
	if !ok || t == "" {
		t, ok = request.Headers["authorization"]
		if !ok || t == "" {
			return "", fmt.Errorf("no authorization header: %w", serrors.ErrAuthorization)
		}
	}

	token, err := h.parser.ParseToken(t)
	if err != nil {
		return "", err
	}

	return token.userID, nil
}

type JwtParser struct{}

type ParsedToken struct {
	userID string
}

func (t JwtParser) ParseToken(token string) (ParsedToken, error) {
	s := strings.Split(token, " ")
	if len(s) != 2 {
		return ParsedToken{}, fmt.Errorf("bad bearer token: %w", serrors.ErrBadInput)
	}

	bearer := s[1]
	tk, _, err := jwt.NewParser().ParseUnverified(bearer, &jwt.RegisteredClaims{})
	if err != nil {
		return ParsedToken{}, fmt.Errorf("parse token: %w", err)
	}

	//if !tk.Valid {
	//	return ParsedToken{}, fmt.Errorf("not valid token")
	//}

	cl, ok := tk.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return ParsedToken{}, fmt.Errorf("parse claims: %w", err)
	}

	return ParsedToken{userID: cl.Subject}, nil
}

func ExtractFromPath(request events.APIGatewayV2HTTPRequest, key string) (string, bool) {
	p, ok := request.PathParameters[key]
	return p, ok
}

func ExtractFromQuery(request events.APIGatewayV2HTTPRequest, key string) (string, bool) {
	p, ok := request.QueryStringParameters[key]
	return p, ok
}

func ExtractSliceFromQuery(request events.APIGatewayV2HTTPRequest, key string) ([]string, bool) {
	s, ok := ExtractFromQuery(request, key)
	if !ok {
		return nil, false
	}

	slice := strings.Split(s, ",")
	res := make([]string, 0, len(slice))
	for _, w := range slice {
		res = append(res, strings.TrimSpace(w))
	}

	return res, true
}

func initDB(ctx context.Context) (sqlc.DB, error) {
	dsn := os.Getenv(DsnEnv)
	if dsn == "" {
		return nil, fmt.Errorf("no env param %s", DsnEnv)
	}

	pool, err := pgx2.NewPool(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("init pgx pool: %w", err)
	}

	return pgx2.NewDB(pool, DBSchema), nil
}

func InitLambdaPortfolioHandler(ctx context.Context) (PortfolioHandler, error) {
	db, err := initDB(ctx)
	if err != nil {
		return PortfolioHandler{}, err
	}

	r, err := factory.NewDBCurrencyRater(db)
	if err != nil {
		return PortfolioHandler{}, err
	}

	uc := usecases.NewPortfolioInteractor(
		db,
		postgres.NewPortfolioRepository(),
		r,
		postgres.NewTickerRepository(),
	)

	return NewPortfolioHandler(
		uc,
		//db,
		//postgres.NewPortfolioRepository(),
		//postgres.NewTickerRepository(),
		UserIDExtractorFromAuthHeader{},
	), nil
}

func InitLambdaTickerHandler(ctx context.Context) (TickerHandler, error) {
	db, err := initDB(ctx)
	if err != nil {
		return TickerHandler{}, err
	}

	r, err := factory.NewDBCurrencyRater(db)
	if err != nil {
		return TickerHandler{}, err
	}

	uc := usecases.NewTickerInteractor(
		db,
		postgres.NewTickerRepository(),
		r,
	)

	return NewTickerHandler(
		uc,
		//db,
		//postgres.NewTickerRepository(),
	), nil
}
