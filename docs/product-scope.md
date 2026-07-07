# Product Scope

## Goal

Build a very small marketplace MVP in controlled phases. The current implemented foundation supports seller authentication, seller-owned product management, product publishing, and a public catalog of published products.

## Implemented Now

### Seller Account

- Sellers can register for an account.
- Sellers can log in.

### Seller Product Management

- Sellers can create draft products.
- Sellers can list their own products.
- Sellers can view details for their own products.
- Sellers can publish their own draft products.

### Public Catalog

- Visitors can list published products.
- Visitors can view published product details.
- Draft products are hidden from public catalog routes.

## Deferred Future Scope

- Product images.
- Product variants.
- Categories.
- Search and pagination.
- Buyer accounts and buyer-owned storefront behavior.
- Product unpublish.

## Out of Scope

- Cart
- Orders
- Payments
- Shipping
- Promotions
- Reviews
- Multi-warehouse inventory
- Microservices
- AI features
- Automated deployment, push, or merge behavior

## Current Boundary

The current repository stops at one-way draft-to-published product publishing and a simple public catalog. Product images, variants, categories, buyer accounts, checkout, fulfillment, post-purchase workflows, seller marketing tools, and complex platform architecture require explicit future task scope before implementation.
