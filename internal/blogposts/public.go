package blogposts

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/predicate"
)

// StatePublished is the `state` of a post open to readers. The ent field is a
// plain string, so the value queries compare against is named here, apart from
// the contract's enum.
const StatePublished = "published"

// RecommendedLimit is how many other posts the post page offers (legacy
// `.limit(2)`).
const RecommendedLimit = 2

// Published is legacy `BlogPost.published_state.with_locale`: what a reader may
// see in locale. A draft, an archived post or a post in another language is
// simply not found.
func Published(locale string) predicate.BlogPost {
	return blogpost.And(blogpost.StateEQ(StatePublished), blogpost.LocaleEQ(locale))
}

// Public reads the posts a visitor may open, newest (highest id) first.
type Public struct {
	db *ent.Client
}

// NewPublic builds the public blog reader over an ent client.
func NewPublic(db *ent.Client) *Public {
	return &Public{db: db}
}

// Query is the published posts of locale with their creator, newest first, for a
// caller that pages them.
func (p *Public) Query(locale string) *ent.BlogPostQuery {
	return p.db.BlogPost.Query().
		Where(Published(locale)).
		WithCreator().
		Order(ent.Desc(blogpost.FieldID))
}

// BySlug is the post a reader asked for; a miss is ent not-found.
func (p *Public) BySlug(ctx context.Context, locale, slug string) (*ent.BlogPost, error) {
	return p.Query(locale).Where(blogpost.SlugEQ(slug)).Only(ctx)
}

// Recommended is the posts to read next after post: newest first, never the post
// itself. Divergence from legacy, on purpose: it wrote `.except(blog_post)`,
// which is ActiveRecord's query-part remover, not a record exclusion — a no-op,
// so its two posts were unordered and could include the post being read.
func (p *Public) Recommended(ctx context.Context, locale string, post *ent.BlogPost) ([]*ent.BlogPost, error) {
	return p.Query(locale).Where(blogpost.IDNEQ(post.ID)).Limit(RecommendedLimit).All(ctx)
}

// NextOlder is the next older published post (id strictly below the given one),
// legacy `Api::BlogPostsController#next`. The oldest post has none: ent
// not-found, which ends the page's infinite scroll.
func (p *Public) NextOlder(ctx context.Context, locale string, id int) (*ent.BlogPost, error) {
	return p.Query(locale).Where(blogpost.IDLT(id)).First(ctx)
}
