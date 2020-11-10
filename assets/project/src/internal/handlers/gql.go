package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/99designs/gqlgen/graphql/handler"
	"{{ . }}/graphql"
	"{{ . }}/graphql/index"
	gql "{{ . }}/internal/gql/generated"
	"{{ . }}/internal/gql/resolvers"
	th "{{ . }}/pkg/transports/http"
	"{{ . }}/pkg/validator"
)

// GraphqlHandler defines the GQLGen GraphQL server handler
func GraphqlHandler(db *gorm.DB, v *validator.Validator) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	c := gql.Config{
		Resolvers: &resolvers.Resolver{
			DB: db, // pass in the ORM instance in the resolvers to be used
			V:  v,
		},
	}

	h := handler.NewDefaultServer(gql.NewExecutableSchema(c))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

var page = template.Must(vfstemplate.ParseFiles(index.Index, nil, "graphiql.tmpl"))

// var page = template.Must(template.New("graphiql").Parse())

func playgroundHandler(title string, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		err := page.Execute(w, map[string]string{
			"title":    title,
			"endpoint": endpoint,
		})
		if err != nil {
			panic(err)
		}
	}
}

// PlaygroundHandler defines a handler to expose the Playground
func PlaygroundHandler(path string) gin.HandlerFunc {
	h := playgroundHandler("Go GraphQL Server", path)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

type Options struct {
	GraphQLPath         string
	IsPlaygroundEnabled bool
	PlaygroundPath      string
}

func CreateGqlHandlers(
	o *Options,
	logger *zap.Logger,
	db *gorm.DB,
	v *validator.Validator,

) th.InitHandlers {
	return func(r *gin.Engine) {
		r.StaticFS("/assets", graphql.Assets)

		apiv1 := r.Group("/api")
		apiv1.POST(o.GraphQLPath, GraphqlHandler(db, v))
		logger.Info("GraphQL @ ", zap.String("graphql path", o.GraphQLPath))
		// Playground handler
		if o.IsPlaygroundEnabled {
			logger.Info("GraphQL Playground @ ", zap.String("graphql path", r.BasePath()+o.PlaygroundPath))
			apiv1.GET(o.PlaygroundPath, PlaygroundHandler(filepath.Join(apiv1.BasePath(), o.GraphQLPath)))
		}
	}
}
