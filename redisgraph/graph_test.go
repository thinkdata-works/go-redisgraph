package redisgraph

import(
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
}

func (suite *Suite) SetupDB() {

}

func (suite *Suite) BeforeTest(suitename string, testname string) {
	suite.SetupDB()
}

func TestGraph(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) TestGraphQuery1() {
	// TODO move all graph stuff to test hooks
	// Open graph connection
	conn, err := redis.Dial("tcp", "redisgraph:6379")
	defer conn.Close()
	suite.Assert().NoError(err)

	// TODO defer cleaing the graph

	// Insert some things
	graph := CreateGraph(&conn, "graph-test")
	res, err := graph.Query("CREATE (:person {name: 'Steve Albini'})")
	suite.Assert().NoError(err)
	fmt.Sprintf("%+v", res)
}