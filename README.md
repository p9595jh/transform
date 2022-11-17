# transform

## Description

`transform` can make a golang object transform its fields or be mapped to another object by tagging.

And also registering more transformers is possible.

It uses generic so Go1.18+ needed.

## Tag

`transform` uses tag to execute transformation commands. Default tag name is `transform` and it can be changed.

The basic shape of tagging is like below:

```
type Person struct {
	Name   string `transform:"upper,bytes"`
	Age    int    `transform:"x2,x:3,big"`
}
```

You can check the tag by `Tag` method, and can change this by `SetTag` method.

```
// declare transformer
transformer := transform.New()

// check tag (default: 'transform')
tag := transformer.Tag()
fmt.Println("tag is", tag)

// change tag
transformer.SetTag("sample")
// then now you only can use this transformer to parse 'sample' tag

// change tag when declare
transformer := transform.New().SetTag("hello")
```

Tag is composed of multiple commands with comma(,).

If a transformation needs a parameter, you can define it using colon(:).

```
type Sample struct {
    A int `transform:"add1"`       // this just runs 'add1' transformer
    B int `transform:"add:1"`      // this runs 'add' transformer and give a parameter '1'
    C int `transform:"add1,add:1"` // multiple tags also allowed (combined with comma)
}
```

## Register

Custome transformers can be registered by `RegisterTransformer` method.
This method receives a name and transformation function.
The transformation function only can be implemented by the function `F1` and `F2`.

`F1` receives one generic type, so the type of input and output is same.
This is used for just transformation.

```
type Sample struct {
    V1 int `transform:"add:10"`
    V2 int `transform:"add:3"`
}

transformer.RegisterTransformer("add", transform.F1[int](func(i int, s string) int {
    x, _ := strconv.ParseInt(s, 10, 32)
    return i + int(x)
}))
```

Above example is registering 'add' transformer. This will update fields as the given parameter.
The second parameter of the function `F1` and `F2` is, as you already knew, the parameter.

The only difference between `F1` and `F2` is a number of generic type. `F2` can be used when mapping, so it needs input type and output type seperately.

```
type Original struct {
    V int `transform:"big"`
}

type Mapped struct {
    V *big.Int
}

transformer.RegisterTransformer("big", transform.F2(func(i int, s string) *big.Int {
	return big.NewInt(int64(i))
}))
```

Transformer also can be registered when initializing.

```
transformer := transform.New(
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
```

## Transform

`Transform` transforms an object. It returns error when there was an error.

```
type Sample struct {
    S string `transform:"upper"`
}

sample := Sample{"hello"}
transformer := transform.New()
if err := transformer.Transform(&sample); err != nil {
    panic(err)
}

fmt.Println(sample) // "HELLO"
```

## Mapping

`mapping` maps `src` to `dst`. It returns error when there was an error.

```
type Original struct {
    V int `transform:"big"`
}

type Mapped struct {
    V *big.Int
}

transformer := transform.New()
transformer.RegisterTransformer("big", transform.F2(func(i int, s string) *big.Int {
	return big.NewInt(int64(i))
}))

original := Original{10}
var mapped Mapped

if err := transformer.Mapping(&original, &mapped); err != nil {
    panic(err)
}

fmt.Println(mapped) // {V: *big.Int(10)}
```

## DTO example

This example can be found in [transform_test.go](./transform_test.go).

```
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
```

## How to use

```
go get github.com/p9595jh/transform
```

```
import "github.com/p9595jh/transform"
```
