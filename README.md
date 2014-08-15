# CODL: Cookoo Domain Language

CODL (pronounced *coddle*) is a simple domain specific language (DSL) for writing Cookoo
routes. It is designed to feel like SQL and read like English.

Like SASS, Thrift, and other similar DSLs, CODL files are transformed
into code (Go code, in this case) prior to compile time.

## Syntax

Here is a basic example of the syntax:

```
IMPORT
	github.com/Masterminds/cookoo/web
	github.com/Masterminds/cookoo/cli

ROUTE "test" "This is a test route."
  DOES web.Flush "first"
    USING p1 "This is a default value"
    USING p2 `1` FROM cxt:p2
  DOES cli.ParseArgs "CMD"

ROUTE "foo" "This is another route"
  DOES web.Flush cmd3
```


## Keywords

CODL provides the following commands. *Case is important!* These MUST be
in all caps.

- `IMPORT`: Import one or more Go packages.
- `ROUTE`: Add a new route
- `DOES`: Add a command to a route
- `USING`: Set a parameter on a command, and optionally set a default
- `FROM`: Pass a value into a parameter on a command
- `INCLUDES`: Include another route in the present route.

CODL cannot tell bare words (see below) from statements. So if you need
to use a string that exactly matches a statement name, make sure you
enclose it in quotation marks.

```
Keyword:
IMPORT

String:
"IMPORT"
```

## Strings

CODL supports three kinds of strings:

1. Plain strings: `"This is a string"`
2. Bare words (unquoted single-word strings): `thisIsAString`
3. Code strings: `«1 + 3»` (or, if you prefer, surround them with
   backticks instead if double angle brackets).

### Plain Strings

A plain string is a quotation mark, followed by any series of Unicode characters
(including line breaks) that is not a quotation mark, followed by a
closing quotation mark.

```
"This is a string"

'this is also a string'

"And it's okay to do this."

"©ƒ†ç˚ˆ˙"

// Backslashes are escapes in double- and single-quotes.
"That\'s it!" // Becomes "That's it!"

"Double\\backslash" // Becomes "Double\backslash"
```

### Bare Words

A bare word is an unquoted string with no whitespace characters.

```
bareword

bare+word

/a/path/is/a/bare/word
```
Bare words are designed as a convenience, and to be *most* convenient,
the CODL interpreter will... well... *interpret* what you mean by them.
Thus, in some cases, a bare word may be assumed to be a code literal,
while in other cases it may be assumed to be a quoted string.

The current set of rules for this is simple:

* If the bare word appears immediately after `DOES`, it is considered a
  code literal: `DOES foo.Bar`
* In all other cases it is assumed to be a string.


### Code Literals

A code literal (or code string) is a sequence of characters that
represents a piece of Go code. It is inserted unaltered into the
generated Go code.

There are two forms of code literals:

Backtick strings:
```
`1` // An int

`true` // a bool

`fmt.Printf("Foo %f", 2.4)` // a function
```

Double-angle quotes:

```
«1» // An int

«true» // a bool

«fmt.Printf("Foo %f", 2.4)» // a function
```

Why do we use double-angles? Because this minimizes your needing to
escape sequences in your code, and also makes it easy to embed CODL
inside of multi-line Go strings.

## Statements

The following statements can be built using keywords and strings:

### IMPORT

```
IMPORT [string [string [string [...]]]]
```

`IMPORT` can occur any number of times, but only at the top of a file:

```
IMPORT foo
IMPORT bar baz

ROUTE //...
```

The above will generate:

```go
import (
  "foo"
  "bar"
  "baz"
)

// ...
```

### ROUTE

`ROUTE` is the main command available in CODL. A route command is
composed of the following pieces. (All but ROUTE are optional)

```
ROUTE "name string" "description or help text as a string"
  DOES `commandAsCodeLiteral` "name string"
    USING "param name" "optional default value (may also be a code literal)"
      FROM "one" "or" "more" "strings"
```

The above will generate the following Go code (slightly altered for
readability).

```go
package routes

import (
  "github.com/Masterminds/cookoo"
)

func Routes(reg *cookoo.Registry) {
  registry.Route("name string", "description or help text as a string").
    Does(commandAsCodeLiteral, "name string").
      Using("param name").WithDefault("optional default value...").
        From("one", "or", "more", "strings")
}
```

### INCLUDES

A `ROUTE` can also include anther route with `INCLUDES`.

```
ROUTE a "First route"

ROUTE b "Second route includes first"
  INCLUDES a

```

Roughly, this would produce:

```go

reg.Route("a", "First Route")

reg.Route("b", "Second route includes first").
  Includes("a")
```

As a general rule of thumb, you should always declare a route before
including it elsewhere (though honestly CODL doesn't care).

## Whitespace

Outside of strings, CODL treats whitespace as significant only as a
separator.

Thus, `ROUTE foo bar DOES foo bar USING foo bar FROM foo bar` is treated as
semantically identical to:

```
ROUTE foo bar
  DOES foo bar
    USING foo bar
      FROM foo bar
```

CODL does not care about your indentation practices. Use tabs. Use
spaces. Use a carefully crafted sequence of carriage returns. And
bikeshed away on which one is best.

## Comments

There is one style of comment in CODL files: `//`. A comment extends
from `//` to the end of the line (`\n`).

Example:

```
// Do imports
IMPORT foo // That's foo.
```

There is one **important** exception for comments. Because slashes
appear frequently in URLs and paths, a comment cannot be immediately
adjacent to a bare word:

```
ROUTE name description//foo
```

In the above, `description//foo` will be treated as one string. This is
very useful for cases like this:

```
USES url http://example.com
```

As a result of this, attempting to start a comment immediately adjacent
to a keyword is likely to cause a parse error:

```
// Bad.
IMPORTS//comment
IMPORTS// comment

// Good.
IMPORTS //comment
IMPORTS // comment
```

