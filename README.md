# go-unicode
parse unicode to string

## example
```go
import (
 unicode "github.com/cocktail18/unicode"
)


unicodeStr := "\\\\ud83d\\\\udc37\\\\ud83d\\\\udc37"
fmt.Println(unicode.UnicodeToString(unicodeStr))
// output 🐷🐷 <nil>
```
