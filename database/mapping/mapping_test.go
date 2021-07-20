package mapping

import (
    "errors"
    "fmt"
    "reflect"
    "testing"

    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

type TestStruct struct {
    IntVal int
    StrVal string
}

type TestStructTags struct {
    IntVal  int64 `dbfield:"i"`
    IntVal2 int8  `dbMapName:"j"`
    IntVal3 int16 `json:"k"`
    IntVal4 int32
}

func TestGetMapping(t *testing.T) {
    m, err := GetMapping(TestStruct{})
    require.Nil(t, err)
    require.Equal(t, reflect.TypeOf(TestStruct{}), m.GetType())
    fieldNames := m.GetFieldNames()
    require.Equal(t, "IntVal", fieldNames[0])
    require.Equal(t, "StrVal", fieldNames[1])
    fields := m.GetFields()
    f, ok := fields["IntVal"]
    require.True(t, ok)
    require.Equal(t, reflect.TypeOf(int(1)), f.Properties.Type)
    f2, ok2 := fields["StrVal"]
    require.True(t, ok2)
    require.Equal(t, reflect.TypeOf(string("")), f2.Properties.Type)

    m2, err := GetMapping(TestStruct{})
    require.Nil(t, err)
    require.Same(t, m, m2)

    // TODO: how to test for non-existence of race condition after mapping lock?

    m, err = GetMapping(&TestStruct{})
    require.Nil(t, err)
    require.Equal(t, reflect.TypeOf(TestStruct{}), m.GetType())
    fieldNames = m.GetFieldNames()
    require.Equal(t, "IntVal", fieldNames[0])
    require.Equal(t, "StrVal", fieldNames[1])
    fields = m.GetFields()
    f, ok = fields["IntVal"]
    require.True(t, ok)
    require.Equal(t, reflect.TypeOf(int(1)), f.Properties.Type)
    f2, ok2 = fields["StrVal"]
    require.True(t, ok2)
    require.Equal(t, reflect.TypeOf(string("")), f2.Properties.Type)

    m, err = GetMapping(TestStructTags{})
    require.Nil(t, err)
    fields = m.GetFields()
    f, ok = fields["i"]
    require.True(t, ok)
    require.Equal(t, reflect.TypeOf(int64(1)), f.Properties.Type)
    f2, ok = fields["j"]
    require.True(t, ok)
    require.Equal(t, reflect.TypeOf(int8(1)), f2.Properties.Type)
    f3, ok := fields["k"]
    require.True(t, ok)
    require.Equal(t, reflect.TypeOf(int16(1)), f3.Properties.Type)
    f4, ok := fields["IntVal4"]
    require.True(t, ok)
    require.Equal(t, reflect.TypeOf(int32(1)), f4.Properties.Type)

    m, err = GetMapping(1)
    require.NotNil(t, err)
    require.Nil(t, m)
}

func TestValuesOf(t *testing.T) {
    v := TestStruct{1, "1"}
    vs, err := ValuesOf(v)
    require.Nil(t, err)
    require.NotNil(t, vs)
    require.Equal(t, 2, len(vs))
    require.IsType(t, int(10), vs[0])
    require.IsType(t, string("10"), vs[1])

    v2 := &TestStruct{2, "2"}
    vs, err = ValuesOf(v2)
    require.Nil(t, err)
    require.NotNil(t, vs)
    require.Equal(t, 2, len(vs))
    require.IsType(t, int(20), vs[0])
    require.IsType(t, string("20"), vs[1])

    vs, err = ValuesOf(1)
    require.NotNil(t, err)
    require.Nil(t, vs)

    // TODO: how to test field.CanInterface()==false ?
}

/*
TestRows

Mocks the functionality of sql.Row for Next() and Scan(...) calls.

Err() panics as it should never be called in this test.
*/
type TestRows struct {
    mock.Mock
}

func (t *TestRows) Next() bool {
    r := t.Called()
    return r.Get(0).(bool)
}
func (t *TestRows) Scan(params ...interface{}) error {
    r := t.Called(params)
    if r.Get(0) != nil {
        return r.Get(0).(error)
    }
    return nil
}
func (t *TestRows) Err() error {
    panic("Err() should not be called")
}

func TestScan(t *testing.T) {
    // Grab mapping
    m, _ := GetMapping(TestStruct{})

    // Re-usable argument matcher
    scanArgMatch := mock.MatchedBy(func([]interface{}) bool { return true })

    // Testing successful Scan
    mockRows := new(TestRows)
    mockRows.On("Scan", scanArgMatch).Return(nil)
    v, err := m.Scan(mockRows)
    require.Nil(t, err)
    require.NotNil(t, v)
    _, ok := v.(*TestStruct)
    require.True(t, ok)
    mockRows.AssertNumberOfCalls(t, "Scan", 1)

    // Testing unsuccessful Scan (no next row)
    testErr := errors.New("testError")
    mockRows = new(TestRows)
    mockRows.On("Scan", scanArgMatch).Return(testErr)
    v, err = m.Scan(mockRows)
    require.Equal(t, testErr, err)
    require.Nil(t, v)
    mockRows.AssertNumberOfCalls(t, "Scan", 1)

    // Testing Multiscan limit=10
    mockRows = new(TestRows)
    mockRows.
        On("Next").Return(true).Times(11).
        On("Scan", scanArgMatch).Return(nil).Times(11)
    vs, err := m.MultiScan(mockRows, 10)
    require.Nil(t, err)
    require.NotNil(t, vs)
    tss, ok := vs.([]*TestStruct)
    require.True(t, ok)
    require.Equal(t, 10, len(tss))
    mockRows.AssertNumberOfCalls(t, "Next", 11)
    mockRows.AssertNumberOfCalls(t, "Scan", 10)

    // Testing Multiscan limit=-1, maxcap=500
    mockRows = new(TestRows)
    i := 0
    scanFn := func(args mock.Arguments) {
        params := args.Get(0).([]interface{})
        *(params[0].(*int)) = i
        *(params[1].(*string)) = fmt.Sprintf("%d", i)
        i++
    }
    mockRows.
        On("Next").Return(true).Times(500).
        On("Next").Return(false).
        On("Scan", scanArgMatch).Run(scanFn).Return(nil)
    vs, err = m.MultiScan(mockRows, -1)
    require.Nil(t, err)
    require.NotNil(t, vs)
    tss, ok = vs.([]*TestStruct)
    require.True(t, ok)
    require.Equal(t, 500, len(tss))
    for j := 0; j < 500; j++ {
        require.Equal(t, j, tss[j].IntVal)
        require.Equal(t, fmt.Sprintf("%d", j), tss[j].StrVal)
    }
    mockRows.AssertNumberOfCalls(t, "Next", 501)
    mockRows.AssertNumberOfCalls(t, "Scan", 500)

    // should fail after 5 because of err
    mockRows = new(TestRows)
    mockRows.
        On("Next").Return(true).
        On("Scan", scanArgMatch).Return(nil).Times(5).
        On("Scan", scanArgMatch).Return(testErr)
    vs, err = m.MultiScan(mockRows, -1)
    require.Equal(t, testErr, err)
    require.NotNil(t, vs)
    tss, ok = vs.([]*TestStruct)
    require.True(t, ok)
    require.Equal(t, 5, len(tss))
    mockRows.AssertNumberOfCalls(t, "Next", 6)
    mockRows.AssertNumberOfCalls(t, "Scan", 6)

    // multi scan with limit=0 should yield empty result and no calls on row
    mockRows = new(TestRows)
    mockRows.
        On("Next").Return(true).
        On("Scan", scanArgMatch).Return(nil)
    vs, err = m.MultiScan(mockRows, 0)
    require.Nil(t, err)
    require.NotNil(t, vs)
    tss, ok = vs.([]*TestStruct)
    require.True(t, ok)
    require.Equal(t, 0, len(tss))
    mockRows.AssertNumberOfCalls(t, "Next", 0)
    mockRows.AssertNumberOfCalls(t, "Scan", 0)
}

/**
 * Test FieldProperties
 * MarshalJSON
 * */

func TestFieldPropertiesMarshalJSON(t *testing.T) {
    fp := FieldProperties{Type: reflect.TypeOf([]int{})}
    bytes, err := fp.MarshalJSON()
    require.Nil(t, err)
    require.Equal(t, `{"Type":"[]int"}`, string(bytes))
}

/**
 * Test helper functions
 *
 * */

func TestMust(t *testing.T) {
    require.Panics(t, func() { must(errors.New("test")) })
    require.NotPanics(t, func() { must(nil) })
}
