# CODL: Cookoo Domain Language


## Syntax

Here is a basic example of the syntax:

```
IMPORT "foo/some"

ROUTE "name" "Description"
  DOES `some.Command` "name"
    USING "paramName" "Optional default" FROM cxt:foo post:bar
```

Syntax:

Statements:

- IMPORT
- ROUTE
- DOES
- USING
- FROM
- INCLUDES

Quotes:

```
"This is a string"

unquoted_string
```

Code Literals

```
`1`

`true`

`fmt.Printf("Foo %f", 2.4)`
```


