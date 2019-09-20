package graph

import (
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
	conn redis.Conn
	graph *Graph
}

func (suite *Suite) setupdb() {
	conn, err := redis.Dial("tcp", "redisgraph:6379")
	suite.Assert().NoError(err)

	suite.conn = conn

	// Create graph connection
	graph := CreateGraph(&conn, "graph-test")

	// Create a dummy node to initialize graph
	_, err = graph.Query("CREATE (:node)")
	suite.Assert().NoError(err)

	_, err = graph.Query("MATCH (n) DELETE n")
	suite.Assert().NoError(err)
	suite.graph = graph
}

func (suite *Suite) BeforeTest(suitename string, testname string) {
	suite.setupdb()
}

func TestGraph(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) TestSimpleCreateQuery() {
	defer suite.conn.Close()
	asserts := suite.Assert()

	res, err := suite.graph.Query("CREATE (:person {name: 'Steve Albini', band: 'Shellac'})")
	asserts.NoError(err)

	// Should be an empty response
	asserts.Equal(0, len(res.Headers))
	asserts.Equal(0, len(res.Rows))

	res, err = suite.graph.Query("CREATE (:person {name: 'Kim Gordon', band: 'Sonic Youth'})")
	asserts.NoError(err)

	// Should be an empty response
	asserts.Equal(0, len(res.Headers))
	asserts.Equal(0, len(res.Rows))

	res, err = suite.graph.Query("CREATE (:person {name: 'Frank Black', band: 'Pixies'})")
	asserts.NoError(err)

	// Should be an empty response
	asserts.Equal(0, len(res.Headers))
	asserts.Equal(0, len(res.Rows))

	res, err = suite.graph.Query("MATCH (p:person) RETURN p.name as name, p.band as band ORDER BY name")
	asserts.NoError(err)

	// Should contain the name
	asserts.Equal(2, len(res.Headers))
	asserts.Equal("name", res.Headers[0])
	asserts.Equal("band", res.Headers[1])
	asserts.Equal(3, len(res.Rows))
	asserts.Equal(2, len(res.Rows[0]))
	suite.assertCellEqualsString(res.Rows[0][0], "Frank Black")
	suite.assertCellEqualsString(res.Rows[0][1], "Pixies")
	suite.assertCellEqualsString(res.Rows[1][0], "Kim Gordon")
	suite.assertCellEqualsString(res.Rows[1][1], "Sonic Youth")
	suite.assertCellEqualsString(res.Rows[2][0], "Steve Albini")
	suite.assertCellEqualsString(res.Rows[2][1], "Shellac")

	// Test getting all cells for values
	names, err := res.CellsFor("name")
	asserts.NoError(err)
	asserts.Equal(3, len(names))
	suite.assertCellEqualsString(*names[0], "Frank Black")
	suite.assertCellEqualsString(*names[1], "Kim Gordon")
	suite.assertCellEqualsString(*names[2], "Steve Albini")

	bands, err := res.CellsFor("band")
	asserts.NoError(err)
	asserts.Equal(3, len(bands))
	suite.assertCellEqualsString(*bands[0], "Pixies")
	suite.assertCellEqualsString(*bands[1], "Sonic Youth")
	suite.assertCellEqualsString(*bands[2], "Shellac")

	_, err = res.CellsFor("instruments")
	if asserts.Error(err) {
		asserts.Equal("Header not found", err.Error())
	}
}

func (suite *Suite) TestCreateQueryAllTypes() {
	defer suite.conn.Close()
	asserts := suite.Assert()

	res, err := suite.graph.Query("CREATE (:person {id: 1, name: 'Stephen Malkmus', bands: ['Pavement', 'Silver Jews'], age: 53, confirmed_mensch: true})")
	asserts.NoError(err)

	// Should be an empty response
	asserts.Equal(0, len(res.Headers))
	asserts.Equal(0, len(res.Rows))

	// Pull all attributes and test type conversions
	res, err = suite.graph.Query("MATCH (p:person) RETURN p.id as id, p.name as name, p.bands as bands, p.age as age, p.confirmed_mensch as is_top_dog")
	asserts.NoError(err)

	// Should contain all headers
	asserts.Equal(5, len(res.Headers))
	asserts.Equal([]string{"id", "name", "bands", "age", "is_top_dog"}, res.Headers)

	// Should contain all results
	asserts.Equal(1, len(res.Rows))
	asserts.Equal(5, len(res.Rows[0]))

	suite.assertCellEqualsInt(res.Rows[0][0], 1)
	suite.assertCellEqualsString(res.Rows[0][1], "Stephen Malkmus")
	suite.assertCellEqualsString(res.Rows[0][2], "[Pavement, Silver Jews]")
	suite.assertCellEqualsInt(res.Rows[0][3], 53)
	suite.assertCellTrue(res.Rows[0][4])
}

func (suite *Suite) TestCreateQueryNulls() {
	defer suite.conn.Close()
	asserts := suite.Assert()

	res, err := suite.graph.Query("CREATE (:person {id: 1, name: '', bands: [], age: 0})")
	asserts.NoError(err)

	// Should be an empty response
	asserts.Equal(0, len(res.Headers))
	asserts.Equal(0, len(res.Rows))

	// Pull all attributes and test type conversions
	res, err = suite.graph.Query("MATCH (p:person) RETURN p.id as id, p.name as name, p.bands as bands, p.age as age, p.confirmed_mensch as is_top_dog")
	asserts.NoError(err)

	// Should contain all headers
	asserts.Equal(5, len(res.Headers))
	asserts.Equal([]string{"id", "name", "bands", "age", "is_top_dog"}, res.Headers)

	suite.assertCellEqualsInt(res.Rows[0][0], 1)
	suite.assertCellEqualsString(res.Rows[0][1], "")
	suite.assertCellEqualsString(res.Rows[0][2], "[]")
	suite.assertCellEqualsInt(res.Rows[0][3], 0)
	asserts.True(res.Rows[0][4].IsNull())
}

func (suite *Suite) assertCellEqualsString(cell ResultCell, value string) {
	str, err := cell.ToString()
	suite.Assert().NoError(err)
	suite.Assert().Equal(value, str)
}

func (suite *Suite) assertCellEqualsInt(cell ResultCell, value int) {
	intv, err := cell.ToInt()
	suite.Assert().NoError(err)
	suite.Assert().Equal(value, intv)
}

func (suite *Suite) assertCellTrue(cell ResultCell) {
	boolv, err := cell.ToBool()
	suite.Assert().NoError(err)
	suite.Assert().True(boolv)
}

func (suite *Suite) assertCellFalse(cell ResultCell) {
	boolv, err := cell.ToBool()
	suite.Assert().NoError(err)
	suite.Assert().False(boolv)
}