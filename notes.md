# Notes
things to look out while using go buffalo

## Testing
1. passing UUID as string (`%s`), in the examples it is shown as `%d` which is wrong

```go
	res := as.HTML("/posts/%s", post.ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), post.Content)
```

2. GO & ORMs
- many of the ORMs don't have all the features like ones found in other dynamic languages
- GORM, Beego's ORM are having most of the features (from go libs that I know so far)
- Pop uses sqlx behind the scenes and doesn't have much support/example for handling relations.
