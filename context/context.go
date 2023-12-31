package context

import (
	"context"
	"fmt"
	"lenslocked/M"
)

type Key string
const UserKey Key="user"
const GalleryKey="gallery"

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

func WithGallery(ctx context.Context ,gallery *M.Gallery)context.Context{
	return context.WithValue(ctx,GalleryKey,gallery)
}

func Gallery(ctx context.Context)(*M.Gallery,error){
	gallery,ok:=ctx.Value(GalleryKey).(*M.Gallery)
	if !ok{
		return nil,fmt.Errorf("gallery not found")
	}
	return gallery,nil
}