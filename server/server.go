package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	graphql "github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	"github.com/rhino1998/umdplanner/testudo"
)

func main() {
	f, err := os.Open("testudo.json")
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	store, err := testudo.LoadStore(f)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	schema := graphql.MustParseSchema(schemaString, &Resolver{store})

	s, err := newServer(store)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	r := chi.NewRouter()

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))
	r.Handle("/query", &relay.Handler{Schema: schema})
	r.Get("/classes/query", s.QueryCoursesGET)

	http.ListenAndServe(":3001", r)

}

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
				return fetch("/query", {
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
