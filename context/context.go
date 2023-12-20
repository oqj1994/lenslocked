package context

import (
	"context"
	"fmt"
	"lenslocked/M"
)

type Key string
const UserKey Key="user"

func WithUser(ctx context.Context ,user *M.User)context.Context{
	return context.WithValue(ctx,UserKey,user)
}

func User(ctx context.Context)(*M.User,error){
	user,ok:=ctx.Value(UserKey).(*M.User)
	if !ok{
		return nil,fmt.Errorf("user not found")
	}
	return user,nil
}