package auth

import "context"

type contextKey string

const sellerIDKey contextKey = "seller_id"

func ContextWithSellerID(ctx context.Context, sellerID string) context.Context {
	return context.WithValue(ctx, sellerIDKey, sellerID)
}

func SellerIDFromContext(ctx context.Context) (string, bool) {
	sellerID, ok := ctx.Value(sellerIDKey).(string)
	return sellerID, ok && sellerID != ""
}
