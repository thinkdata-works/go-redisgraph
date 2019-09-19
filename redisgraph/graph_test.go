package redisgraph

import(
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
	conn redis.Conn
	graph *Graph
}

func (suite *Suite) SetupDB() {
	conn, err := redis.Dial("tcp", "redisgraph:6379")
	suite.Assert().NoError(err)

	suite.conn = conn

	// Create graph connection
	graph := CreateGraph(&conn, "graph-test")
	res, err := graph.Query("MATCH (n) DELETE n")
	fmt.Println(fmt.Sprintf("Response clearing graph: %+v", res))
	suite.Assert().NoError(err)
	suite.graph = graph
}

func (suite *Suite) BeforeTest(suitename string, testname string) {
	suite.SetupDB()
}

func TestGraph(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) TestSimpleCreateQuery() {
	defer suite.conn.Close()
	asserts := suite.Assert()

	res, err := suite.graph.Query("CREATE (:person {name: 'Steve Albini'})")
	asserts.NoError(err)

	// Should be an empty response
	asserts.Equal(0, len(res.Headers))
	asserts.Equal(0, len(res.Rows))

	fmt.Println(fmt.Sprintf("%+v", res))

	res, err = suite.graph.Query("MATCH (p:person) RETURN p.name as name")
	asserts.NoError(err)

	// Should contain the name
	asserts.Equal(1, len(res.Headers))
	asserts.Equal("name", res.Headers[0])
	asserts.Equal(1, len(res.Rows)) // one row
	asserts.Equal(1, len(res.Rows[0])) // one cell
	suite.assertCellEqualsString(res.Rows[0][0], "Steve Albini")

	fmt.Println(fmt.Sprintf("%+v", res))
}

func (suite *Suite) assertCellEqualsString(cell ResultCell, value string) {
	str, err := cell.ToString()
	suite.Assert().NoError(err)
	suite.Assert().Equal(value, str)
}