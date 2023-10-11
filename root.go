package graphql

import (
	"embed"
	"fmt"
	nHttp "net/http"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func getPage(mutationPath string) []byte {
	var page = []byte(`
	<!DOCTYPE html>
	<html>
		<head>
			<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
			<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
		</head>
		<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
			<div id="graphiql" style="height: 100vh;">Loading...</div>
			<script>
				function graphQLFetcher(graphQLParams) {
					return fetch("` + mutationPath + `", {
						method: "post",
						body: JSON.stringify(graphQLParams),
						credentials: "include",
					}).then(function (response) {
						return response.text();
					}).then(function (responseBody) {
						try {
							return JSON.parse(responseBody);
						} catch (error) {
							return responseBody;
						}
					});
				}

				ReactDOM.render(
					React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
					document.getElementById("graphiql")
				);
			</script>
		</body>
	</html>
	`)
	return page
}

func NewGraphQLNetHttpHandlerFunc(mutationPath string) nHttp.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(getPage(mutationPath))
		if err != nil {
			panic(fmt.Sprintf("write page error: %s\n", err))
		}
	}
}

func NewNetHttpHandler(root any, content embed.FS) *relay.Handler {
	s, err := String(content)
	if err != nil {
		panic(fmt.Sprintf("reading embedded schema contents: %v", err))
	}

	return &relay.Handler{Schema: graphql.MustParseSchema(s, root, graphql.UseFieldResolvers())}
}

func NewGraphQKratosHttpLHandlerFunc(mutationPath string) http.HandlerFunc {
	return func(ctx http.Context) error {
		_, err := ctx.Response().Write(getPage(mutationPath))
		if err != nil {
			return err
		}
		return nil
	}
}

func NewKratosHttpHandler(root any, content embed.FS) http.HandlerFunc {
	return func(ctx http.Context) error {
		NewNetHttpHandler(root, content).ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}
