package objecthelper

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestSampleTypeJSON struct {
	A string `json:"a"`
}

type TestSampleTypeXML struct {
	A string `xml:"a"`
}

type TestSampleTypeYML struct {
	A string
	B struct {
		C int
	}
}

func TestGetJSONData(t *testing.T) {
	r := strings.NewReader("{\"a\": \"naruto\"}")
	val := TestSampleTypeJSON{}
	err := GetJSONData(ioutil.NopCloser(r), &val)
	assert.NoError(t, err)
	assert.Equal(t, "naruto", val.A)
}

func TestGetDataAsString(t *testing.T) {
	r := strings.NewReader("naruto")
	s, err := GetDataAsString(ioutil.NopCloser(r))
	assert.NoError(t, err)
	assert.Equal(t, "naruto", s)
}

func TestGetJSONDataError(t *testing.T) {
	r := strings.NewReader("{\"a\": \"naruto")
	val := TestSampleTypeJSON{}
	err := GetJSONData(ioutil.NopCloser(r), &val)
	assert.Error(t, err)
}

func TestMarshalUnmarshalJSON(t *testing.T) {
	rq := TestSampleTypeJSON{A: "naruto"}
	var rs TestSampleTypeJSON

	out, err := MarshalJSON(rq)
	assert.NoError(t, err)

	err = UnmarshalJSON(out, &rs)
	assert.NoError(t, err)

	assert.Equal(t, rq, rs)

}

func TestMarshalUnmarshalXML(t *testing.T) {
	rq := TestSampleTypeXML{A: "naruto"}
	var rs TestSampleTypeXML

	out, err := MarshalXML(rq)
	assert.NoError(t, err)

	err = UnmarshalXML(out, &rs)
	assert.NoError(t, err)

	assert.Equal(t, rq, rs)

}

func TestGenerateUUID(t *testing.T) {
	actual, err := GenerateUUID()
	expected, expectedErr := GenerateUUID()
	assert.NotEqual(t, actual, expected)
	assert.Equal(t, err, expectedErr)
}

func TestGetRandomNumber(t *testing.T) {
	n1, err := GetRandomNumber(6)
	assert.NoError(t, err)
	n2, err := GetRandomNumber(6)
	assert.NoError(t, err)
	assert.NotEqual(t, n1, n2)
}

func TestGetXMLData(t *testing.T) {
	s := strings.NewReader("<TestSampleTypeXML><a>naruto</a></TestSampleTypeXML>")
	var rs TestSampleTypeXML
	err := GetXMLData(ioutil.NopCloser(s), &rs)
	assert.NoError(t, err)
	assert.Equal(t, "naruto", rs.A)
}

func TestGetVariableFromContext(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "naruto", "rocks")
	ctx = context.WithValue(ctx, "boruto", true)
	ctx = context.WithValue(ctx, "hinata", "")
	s, err := GetStringVariableFromContext(ctx, "naruto")
	assert.NoError(t, err)
	assert.Equal(t, "rocks", s)
	assert.Equal(t, "rocks", GetStringVariableFromContextD(ctx, "naruto"))
	_, err = GetStringVariableFromContext(ctx, "hinata")
	assert.Error(t, err)
	assert.Equal(t, "missing variable value for hinata", err.Error())
	_, err = GetStringVariableFromContext(ctx, "boruto")
	assert.Error(t, err)
	assert.Equal(t, "invalid variable value for boruto, found true", err.Error())
	_, err = GetStringVariableFromContext(ctx, "sasuke")
	assert.Error(t, err)
	assert.Equal(t, "missing variable sasuke", err.Error())
	b, err := GetBoolVariableFromContext(ctx, "boruto")
	assert.NoError(t, err)
	assert.True(t, b)
	assert.True(t, GetBoolVariableFromContextD(ctx, "boruto"))
	_, err = GetBoolVariableFromContext(ctx, "naruto")
	assert.Error(t, err)
	assert.Equal(t, "invalid variable value for naruto, found rocks", err.Error())
	_, err = GetBoolVariableFromContext(ctx, "sasuke")
	assert.Error(t, err)
	assert.Equal(t, "missing variable sasuke", err.Error())
}

func TestGetYMLData(t *testing.T) {
	s := strings.NewReader(`
a: Easy!
b:
  c: 2
`)
	var rs TestSampleTypeYML
	err := GetYMLData(ioutil.NopCloser(s), &rs)
	assert.NoError(t, err)
	assert.Equal(t, "Easy!", rs.A)
	assert.Equal(t, 2, rs.B.C)
}

func TestMarshalUnmarshalYML(t *testing.T) {
	rs := TestSampleTypeYML{
		A: "naruto",
		B: struct {
			C int
		}{
			C: 8,
		},
	}
	b, err := MarshalYML(rs)
	assert.NoError(t, err)
	var s TestSampleTypeYML
	err = UnmarshalYML(b, &s)
	assert.NoError(t, err)
	assert.Equal(t, rs.A, s.A)
	assert.Equal(t, rs.B.C, s.B.C)
}
