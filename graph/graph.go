package graph

import "github.com/gomodule/redigo/redis"

type Graph struct {
	Name string
	conn *redis.Conn
}

func CreateGraph(conn *redis.Conn, name string) *Graph {
	return &Graph{Name: name, conn: conn}
}

func (g *Graph) Query(query string) (*QueryResult, error) {
	conn := *g.conn
	results, err := conn.Do("GRAPH.QUERY", g.Name, query)
	if err != nil {
		return nil, err
	}

	queryresults, err := createQueryResult(results)
	if err != nil {
		return nil, err
	}

	return queryresults, nil
}
