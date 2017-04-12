# Notes
things to look out while using go buffalo

## Testing
1. passing UUID as string (`%s`), in the examples it is shown as `%d` which is wrong

```go
	res := as.HTML("/posts/%s", post.ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), post.Content)
```

2. 