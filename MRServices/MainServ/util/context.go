package util

import (
	"context"
	"time"

	grpc_tags "github.com/oldjon/gx/modules/grpc/tags"
	"github.com/oldjon/gx/service"
	opentracing "github.com/opentracing/opentracing-go"
)

// PassingContextValues extract logger and tags from an existing ctx and pass them into newCtx
func PassingContextValues(ctxWithValues context.Context, newCtx context.Context) context.Context {
	tags := grpc_tags.FromContext(ctxWithValues)
	newTags := cloneTags(tags)
	newCtx = grpc_tags.WithTags(newCtx, newTags)
	newCtx = opentracing.ContextWithSpan(newCtx, opentracing.SpanFromContext(ctxWithValues))
	return newCtx
}

func CloneContextWithTimeout(ctxWithValues context.Context, dur time.Duration) (context.Context, func()) {
	newCtx, cancel := context.WithTimeout(context.Background(), dur)
	newCtx = PassingContextValues(ctxWithValues, newCtx)
	return newCtx, cancel
}

// ContextWithValues set usingVersion and logger to an existing context
func ContextWithValues(ctx context.Context, logger *service.EventLogger, clientVersion float64, ipRegion string) context.Context {
	tags := grpc_tags.FromContext(ctx)
	newTags := cloneTags(tags).Set(TagsKeyClientVersion, clientVersion)
	newTags = newTags.Set(TagsKeyIPRegion, ipRegion)
	ctx = grpc_tags.WithTags(ctx, newTags)
	return ctx
}

// ContextWithToken set token to an existing context
func ContextWithToken(ctx context.Context, token string) context.Context {
	tags := grpc_tags.FromContext(ctx)
	newTags := cloneTags(tags)
	newTags.Set(TagsKeyToken, token)
	ctx = grpc_tags.WithTags(ctx, newTags)
	return ctx
}

// ContextWithUserTagging set userTagging to an existing context
func ContextWithUserTagging(ctx context.Context, taggingName string, tagging string) context.Context {
	userTaggingFromCtx, ok := ValueFromContext(ctx, TagsKeyUserTagging).([]UserTagging)
	if !ok {
		userTaggingFromCtx = []UserTagging{}
	}

	have := false
	for i := range userTaggingFromCtx {
		if userTaggingFromCtx[i].TaggingName == taggingName {
			userTaggingFromCtx[i].Tagging = tagging
			have = true
			break
		}
	}
	if !have {
		userTag := UserTagging{
			TaggingName: taggingName,
			Tagging:     tagging,
		}
		userTaggingFromCtx = append(userTaggingFromCtx, userTag)
	}

	tags := grpc_tags.FromContext(ctx)
	newTags := cloneTags(tags).Set(TagsKeyUserTagging, userTaggingFromCtx)
	ctx = grpc_tags.WithTags(ctx, newTags)
	return ctx
}

func LoggerWithValues(logger *service.EventLogger, clientVersion float64, ipRegion string) *service.EventLogger {
	tags := grpc_tags.FromContext(context.Background()) // no need to clone because this is already a new Tags *context.Background*
	tags = tags.Set(TagsKeyClientVersion, clientVersion)
	tags = tags.Set(TagsKeyIPRegion, ipRegion)
	newLogger := logger.Clone().SetTags(tags)
	return newLogger
}

func ValueFromContext(ctx context.Context, key string) interface{} {
	tags := grpc_tags.FromContext(ctx)
	return tags.Get(key)
}

// cloneTags: since the context.Context is immutable, passing values from one ctx to another should use the copied values
func cloneTags(tags grpc_tags.Tags) grpc_tags.Tags {
	newTags := grpc_tags.FromContext(context.TODO())
	_ = tags.Foreach(func(key string, val interface{}) error {
		newTags.Set(key, val)
		return nil
	})
	return newTags
}
