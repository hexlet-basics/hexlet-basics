package blog

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/activestorageattachment"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/blogpostlike"
	"hexletbasics/ent/blogpostrelatedcourseitem"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/htmltext"
)

const (
	// wordsPerMinute mirrors the legacy reading-time divisor (BlogPostResource).
	wordsPerMinute = 260
	// recordType/coverName key the polymorphic ActiveStorage rows to a BlogPost
	// (Rails stores the model class name).
	recordType = "BlogPost"
	coverName  = "cover"
)

// Post is a blog post as a reader sees it: the row plus what is not a column on
// it (cover, likes, promoted courses) and what is derived from it (reading
// time). It is the one shape the admin list and the public pages share; the
// handler converts it to the contract and adds nothing.
type Post struct {
	ID          int
	Creator     *ent.User
	Name        *string
	Slug        *string
	Description *string
	State       *string
	Locale      *string
	RichBody    string
	CreatedAt   time.Time

	// ReadingTime is minutes, rounded up; an empty body reads in 0.
	ReadingTime int
	LikesCount  int
	// RelatedCourseIDs are the promoted courses in display order, never nil.
	RelatedCourseIDs []int
	// CoverURL is empty for a post without a cover.
	CoverURL string
}

// Reader assembles posts for every caller that shows them.
type Reader struct {
	db        *ent.Client
	publicURL string
}

// NewReader builds the reader; publicURL is the origin that serves the covers.
func NewReader(db *ent.Client, publicURL string) *Reader {
	return &Reader{db: db, publicURL: publicURL}
}

// Get is one post by id with its creator; a miss is ent not-found.
func (r *Reader) Get(ctx context.Context, id int) (Post, error) {
	row, err := r.db.BlogPost.Query().Where(blogpost.IDEQ(id)).WithCreator().Only(ctx)
	if err != nil {
		return Post{}, errs.Wrapf(err, "load blog post")
	}

	posts, err := r.Assemble(ctx, []*ent.BlogPost{row})
	if err != nil {
		return Post{}, err
	}

	return posts[0], nil
}

// List is one page of posts for the admin, newest first, with the total. A
// non-blank search keeps the posts whose name, slug or description contain it,
// case-insensitively; the body is left out, being the one large text column.
func (r *Reader) List(ctx context.Context, search string, page domain.Page) ([]Post, int, error) {
	query := r.db.BlogPost.Query().WithCreator().Order(ent.Desc(blogpost.FieldID))
	if term := strings.TrimSpace(search); term != "" {
		query = query.Where(blogpost.Or(
			blogpost.NameContainsFold(term),
			blogpost.SlugContainsFold(term),
			blogpost.DescriptionContainsFold(term),
		))
	}

	rows, total, err := domain.Paged(ctx, query, page)
	if err != nil {
		return nil, 0, err
	}

	posts, err := r.Assemble(ctx, rows)

	return posts, total, err
}

// Assemble builds the posts of rows (loaded with their creator), in the order
// given. The cover, like count and promoted courses are not columns on
// blog_posts, so they are fetched in three queries keyed by post id rather
// than per post, keeping a list off an N+1 path.
func (r *Reader) Assemble(ctx context.Context, rows []*ent.BlogPost) ([]Post, error) {
	if len(rows) == 0 {
		return []Post{}, nil
	}

	ids := lo.Map(rows, func(p *ent.BlogPost, _ int) int { return p.ID })

	covers, err := r.coverKeys(ctx, ids)
	if err != nil {
		return nil, err
	}

	likes, err := r.likeCounts(ctx, ids)
	if err != nil {
		return nil, err
	}

	related, err := r.relatedCourseIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	posts := make([]Post, len(rows))
	for i, row := range rows {
		posts[i] = Post{
			ID:               row.ID,
			Creator:          row.Edges.Creator,
			Name:             row.Name,
			Slug:             row.Slug,
			Description:      row.Description,
			State:            row.State,
			Locale:           row.Locale,
			RichBody:         row.RichBody,
			CreatedAt:        row.CreatedAt,
			ReadingTime:      readingTime(row.RichBody),
			LikesCount:       likes[row.ID],
			RelatedCourseIDs: related[row.ID],
		}
		if key := covers[row.ID]; key != "" {
			posts[i].CoverURL = assetstore.URL(r.publicURL, key)
		}
	}

	return posts, nil
}

// coverKeys returns each post's cover blob storage key by id (posts without a
// cover are absent from the map).
func (r *Reader) coverKeys(ctx context.Context, ids []int) (map[int]string, error) {
	rows, err := r.db.ActiveStorageAttachment.Query().
		Where(
			activestorageattachment.RecordType(recordType),
			activestorageattachment.Name(coverName),
			activestorageattachment.RecordIDIn(ids...),
		).
		WithBlob().
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list active storage attachments")
	}

	withBlob := lo.Filter(rows, func(a *ent.ActiveStorageAttachment, _ int) bool { return a.Edges.Blob != nil })

	return lo.SliceToMap(withBlob, func(a *ent.ActiveStorageAttachment) (int, string) {
		return a.RecordID, a.Edges.Blob.Key
	}), nil
}

// relatedCourseIDs returns each post's promoted-course ids in display order
// (the `order` column, id as the tie-breaker for legacy rows without it). Posts
// without items map to an empty slice, so the contract field is always an
// array, never null.
func (r *Reader) relatedCourseIDs(ctx context.Context, ids []int) (map[int][]int, error) {
	rows, err := r.db.BlogPostRelatedCourseItem.Query().
		Where(blogpostrelatedcourseitem.BlogPostIDIn(ids...)).
		Order(
			ent.Asc(blogpostrelatedcourseitem.FieldOrder),
			ent.Asc(blogpostrelatedcourseitem.FieldID),
		).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list blog post related course items")
	}

	out := make(map[int][]int, len(ids))
	for _, id := range ids {
		out[id] = []int{}
	}

	for _, row := range rows {
		out[row.BlogPostID] = append(out[row.BlogPostID], row.CourseID)
	}

	return out, nil
}

// likeCounts returns each post's like count by id via one grouped aggregate.
func (r *Reader) likeCounts(ctx context.Context, ids []int) (map[int]int, error) {
	type row struct {
		BlogPostID int `json:"blog_post_id"`
		Count      int `json:"count"`
	}

	var rows []row

	err := r.db.BlogPostLike.Query().
		Where(blogpostlike.BlogPostIDIn(ids...)).
		GroupBy(blogpostlike.FieldBlogPostID).
		Aggregate(ent.As(ent.Count(), "count")).
		Scan(ctx, &rows)
	if err != nil {
		return nil, errs.Wrapf(err, "read blog post likes")
	}

	return lo.SliceToMap(rows, func(r row) (int, int) { return r.BlogPostID, r.Count }), nil
}

// readingTime is the minutes to read a rich body, rounded up so a short post
// reads in 1 rather than 0 (legacy divided integers and floored; that was
// wrong). An empty body stays 0.
func readingTime(richBodyHTML string) int {
	words := len(strings.Fields(htmltext.PlainText(richBodyHTML)))

	return (words + wordsPerMinute - 1) / wordsPerMinute
}

// Liker records likes.
type Liker struct {
	db *ent.Client
}

// NewLiker builds the like operation over an ent client.
func NewLiker(db *ent.Client) *Liker {
	return &Liker{db: db}
}

// Like records userID's like on the post; a missing post is ent not-found. A
// repeat is not an error and changes nothing: the insert ignores a conflict on
// the unique (post, user) index, so two simultaneous clicks leave one row
// where a check followed by an insert would leave two.
func (l *Liker) Like(ctx context.Context, postID, userID int) error {
	if _, err := l.db.BlogPost.Query().Where(blogpost.IDEQ(postID)).OnlyID(ctx); err != nil {
		return errs.Wrapf(err, "load blog post")
	}

	err := l.db.BlogPostLike.Create().
		SetBlogPostID(postID).
		SetUserID(userID).
		OnConflictColumns(blogpostlike.FieldBlogPostID, blogpostlike.FieldUserID).
		DoNothing().
		Exec(ctx)
	// ent reads the inserted row back, and a skipped insert returns none.
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return oops.Wrapf(err, "like post %d by user %d", postID, userID)
	}

	return nil
}
