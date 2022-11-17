package transform_test

import (
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/p9595jh/transform"
)

type Inner struct {
	What int `transform:"x2,big"`
}

type Person struct {
	Name   string `transform:"upper,bytes"`
	Age    int    `transform:"x2,x:3,big"`
	Nested Inner
}

type InnerDst struct {
	What *big.Int
}

type PersonDst struct {
	Name   []byte
	Age    *big.Int
	Nested InnerDst
}

func TestTransform(t *testing.T) {
	person := Person{"mark", 28, Inner{11}}
	a := transform.New(
		transform.I{
			Name: "x2",
			F: transform.F1(func(i int, s string) int {
				return i * 2
			}),
		},
		transform.I{
			Name: "x",
			F: transform.F1(func(i int, s string) int {
				x, _ := strconv.ParseInt(s, 10, 32)
				return i * int(x)
			}),
		},
	)
	err := a.Transform(&person)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(person)
	}
}

func TestMapping(t *testing.T) {
	person := Person{"mark", 28, Inner{11}}
	a := transform.New()
	a.RegisterTransformer("bytes", transform.F2(func(s1, s2 string) []byte {
		return []byte(s1)
	}))
	a.RegisterTransformer("big", transform.F2(func(i int, s string) *big.Int {
		return big.NewInt(int64(i))
	}))

	var personDst PersonDst
	t.Log(personDst)
	err := a.Mapping(&person, &personDst)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(person, personDst)
	}
}

func TestDto(t *testing.T) {
	// dtom means DTO (data transfer object) Mapper
	type TransactionDtom struct {
		Sender string `json:"sender" transform:"trim0x,lower"`
		Amount string `json:"amount" transform:"big"`
	}

	// TransactionDtom will be mapped to this
	type TransactionDto struct {
		Sender string
		Amount *big.Int
	}

	a := transform.New()
	a.RegisterTransformer("trim0x", transform.F2(func(s1, s2 string) string {
		s1 = strings.TrimPrefix(s1, "0x")
		s1 = strings.ToLower(s1)
		return s1
	}))
	a.RegisterTransformer("big", transform.F2(func(s1, s2 string) *big.Int {
		i := new(big.Int)
		i.SetString(s1, 10)
		return i
	}))

	// raw data -> [unmarshal] -> dtom -> [transform] -> dto
	tx := `{
		"sender":"0x4d943a7C1f2AF858BfEe8aB499fbE76B1D046eC7",
		"amount":"436799733113079832970000"
	}`

	var transactionDtom TransactionDtom
	err := json.Unmarshal([]byte(tx), &transactionDtom)
	if err != nil {
		panic(err)
	}

	var transactionDto TransactionDto
	err = a.Mapping(&transactionDtom, &transactionDto)
	if err != nil {
		panic(err)
	}

	// {sender: 4d943a7c1f2af858bfee8ab499fbe76b1d046ec7, amount: 436799733113079832970000}
	t.Log(transactionDto)
}

func TestTagChange(t *testing.T) {
	type Test struct {
		A string `ttt:"upper"`
	}
	test := Test{"hElLo"}
	a := transform.New().SetTag("ttt")
	t.Log(a.Tag())

	if err := a.Transform(&test); err != nil {
		t.Log(err)
	} else {
		t.Log(test)
	}
}
