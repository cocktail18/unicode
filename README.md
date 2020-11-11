# go-unicode
parse unicode to string

## example
```go
import (
 unicode "cocktail1024.top/go-unicode"
)


unicodeStr := "ab\\\\u5f20\\\\ud83d\\\\udc37\\\\ud83d\\\\udc37c"
fmt.Println(unicode.UnicodeToString(unicodeStr))
// output ab张🐷🐷c nil
```
