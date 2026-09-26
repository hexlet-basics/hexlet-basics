// Package feeds builds the machine-read documents other services poll, whose
// shape is set by the consumer rather than by the contract's JSON models.
package feeds

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/coursecategory"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/coursemoduletranslation"
	"hexletbasics/ent/coursemoduleversion"
	"hexletbasics/ent/landingpage"
)

// feedLocale is the only locale the Yandex feed speaks: Yandex's course
// results are Russian, and the legacy action forced `ru` whatever the URL
// prefix. That is also why the shop strings below are constants rather than
// catalog lookups — they never vary with the request.
const feedLocale = "ru"

// Shop details, as the legacy `common.organization.*` ru strings.
const (
	shopName        = "Code Basics"
	shopCompany     = "ООО «Хекслет Рус»"
	shopEmail       = "support@hexlet.io"
	shopDescription = "Автоматизированная система для обучения программированию в браузере. " +
		"Текстовые курсы + удобный тренажер."
	// The legacy builder hardcoded this absolute URL rather than deriving it
	// from the host, so it stays fixed here too.
	shopPicture = "https://code-basics.com/images/logo.png"
)

// Offer constants the legacy builder hardcoded for every course.
const (
	// yandexCategoryID is the course's rubric in Yandex's education
	// rubricator; every course is filed under the same one.
	yandexCategoryID = 101
	currencyID       = "RUR"
	// minModules is the smallest course Yandex is offered: a course with fewer
	// ru modules is left out of the feed.
	minModules = 3
	// planHours is the duration Yandex shows for each plan module.
	planHours = "5"
	// planIndent reproduces the legacy CDATA byte for byte: the text came from
	// an indented heredoc, so every plan description opens with 20 spaces and
	// ends with a newline.
	planIndent = "                    "
)

// dateLayout is the `yml_catalog date` format Yandex expects.
const dateLayout = "2006-01-02 15:04"

// xmlDeclaration is written by hand: xml.Header adds an encoding attribute
// legacy never sent, and UTF-8 is the XML default either way.
const xmlDeclaration = `<?xml version="1.0"?>` + "\n"

// fixedParams are the course facts the legacy builder asserted for every
// offer, in its order.
var fixedParams = []param{
	{Name: "Продолжительность", Unit: "месяц", Text: "1"},
	{Name: "Формат обучения", Text: "Самостоятельно"},
	{Name: "Сложность", Text: "Для новичков"},
	{Name: "Тип обучения", Text: "Курс"},
	{Name: "Результат обучения", Text: "Сертификат"},
	{Name: "Есть бесплатная часть", Text: "true"},
	{Name: "Есть текстовые уроки", Text: "true"},
	{Name: "Есть тренажеры", Text: "true"},
	{Name: "Есть сообщество", Text: "true"},
}

// Yandex builds the YML course catalogue Yandex reads to list the courses in
// its course results.
type Yandex struct {
	db *ent.Client
	// origin is the canonical site origin the feed's absolute URLs start with.
	origin string
	now    func() time.Time
}

// Option configures a Yandex builder.
type Option func(*Yandex)

// WithClock replaces the wall clock that stamps the catalogue's date, so a
// test can compare the whole document against a golden file.
func WithClock(now func() time.Time) Option {
	return func(y *Yandex) { y.now = now }
}

// NewYandex builds the feed over db, with URLs on appHost — the canonical host
// (config.AppHost), as legacy used AppHost.canonical.
func NewYandex(db *ent.Client, appHost string, opts ...Option) *Yandex {
	y := &Yandex{db: db, origin: "https://" + appHost, now: time.Now}
	for _, opt := range opts {
		opt(y)
	}
	return y
}

// Build renders the catalogue. It is built on every request, as legacy did:
// Yandex polls it rarely, and a cache would only add a staleness window.
func (y *Yandex) Build(ctx context.Context) ([]byte, error) {
	categories, err := y.categories(ctx)
	if err != nil {
		return nil, err
	}
	items, err := y.offers(ctx)
	if err != nil {
		return nil, err
	}

	doc := catalog{
		// Legacy stamped Time.zone.now in the app zone, which is UTC.
		Date: y.now().UTC().Format(dateLayout),
		Shop: shop{
			Name:        shopName,
			Company:     shopCompany,
			URL:         y.origin + "/" + feedLocale,
			Email:       shopEmail,
			Description: shopDescription,
			Picture:     shopPicture,
			Sets:        sets{Set: categories},
			Offers:      offers{Offer: items},
		},
	}

	var buf bytes.Buffer
	buf.WriteString(xmlDeclaration)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return nil, fmt.Errorf("encode yandex feed: %w", err)
	}
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

// categories are the ru course categories, one `set` each. Legacy had no
// ORDER BY and so emitted them in heap order; id order is the deterministic
// stand-in, and Yandex matches sets by id, not position.
func (y *Yandex) categories(ctx context.Context) ([]set, error) {
	rows, err := y.db.CourseCategory.Query().
		Where(coursecategory.LocaleEQ(feedLocale)).
		Order(coursecategory.ByID()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load feed categories: %w", err)
	}
	return lo.Map(rows, func(c *ent.CourseCategory, _ int) set {
		return set{
			ID:   c.ID,
			Name: lo.FromPtr(c.Name),
			URL:  y.localized("/language_categories/" + url.PathEscape(lo.FromPtr(c.Slug))),
		}
	}), nil
}

// offers is one offer per published, listed, main ru landing page of a
// completed course — the legacy `LandingPage.web.where(listed:, main:)` —
// skipping courses whose current version has fewer than minModules ru
// modules.
//
// Legacy had no ORDER BY here either; landing page id is the deterministic
// stand-in.
func (y *Yandex) offers(ctx context.Context) ([]offer, error) {
	pages, err := y.db.LandingPage.Query().
		Where(
			landingpage.LocaleEQ(feedLocale),
			landingpage.StateEQ("published"),
			landingpage.Listed(true),
			landingpage.Main(true),
			landingpage.HasCourseWith(course.ProgressEQ("completed")),
		).
		WithCourse().
		Order(landingpage.ByID()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load feed landing pages: %w", err)
	}

	modules, err := y.modules(ctx, pages)
	if err != nil {
		return nil, err
	}

	var result []offer
	for _, page := range pages {
		crs := page.Edges.Course
		if crs.CurrentVersionID == nil {
			continue
		}
		plan := modules[*crs.CurrentVersionID]
		if len(plan) < minModules {
			continue
		}
		result = append(result, y.offer(page, crs, plan))
	}
	return result, nil
}

// modules loads the ru module translations of every offered course's current
// version in one query, each with its module version's lesson versions and
// their lessons, grouped by course version.
//
// Legacy read `current_module_infos` without an ORDER BY, which in practice is
// insertion order — and the loader inserts modules in build order. Ordering by
// the module version's `order` says that outright; translation id breaks ties.
func (y *Yandex) modules(ctx context.Context, pages []*ent.LandingPage) (map[int][]*ent.CourseModuleTranslation, error) {
	versionIDs := lo.Uniq(lo.FilterMap(pages, func(p *ent.LandingPage, _ int) (int, bool) {
		id := p.Edges.Course.CurrentVersionID
		return lo.FromPtr(id), id != nil
	}))
	if len(versionIDs) == 0 {
		return map[int][]*ent.CourseModuleTranslation{}, nil
	}

	rows, err := y.db.CourseModuleTranslation.Query().
		Where(
			coursemoduletranslation.LocaleEQ(feedLocale),
			coursemoduletranslation.CourseVersionIDIn(versionIDs...),
		).
		WithVersion(func(q *ent.CourseModuleVersionQuery) {
			q.WithLessonVersions(func(q *ent.CourseLessonVersionQuery) {
				q.WithLesson().Order(
					courselessonversion.ByNaturalOrder(sql.OrderNullsLast()),
					courselessonversion.ByID(),
				)
			})
		}).
		Order(
			coursemoduletranslation.ByVersionField(coursemoduleversion.FieldOrder, sql.OrderNullsLast()),
			// Qualify the id: the edge ordering joins language_module_versions,
			// so a bare `id` is ambiguous.
			func(s *sql.Selector) { s.OrderBy(s.C(coursemoduletranslation.FieldID)) },
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load feed modules: %w", err)
	}
	return lo.GroupBy(rows, func(m *ent.CourseModuleTranslation) int { return m.CourseVersionID }), nil
}

// offer renders one course. The offer id is the course id, the landing URL
// uses the landing page slug and the lesson links the course slug, all as
// legacy did.
//
// Legacy added a `picture` for a course with a cover. Course covers are not
// wired in the Go stack yet, so no offer carries one; it comes back with the
// covers.
func (y *Yandex) offer(page *ent.LandingPage, crs *ent.Course, plan []*ent.CourseModuleTranslation) offer {
	slug := lo.FromPtr(page.Slug)
	utm := url.Values{
		"utm_source":   {"yandex"},
		"utm_medium":   {"organic"},
		"utm_content":  {"feed_search"},
		"utm_campaign": {"cb-" + slug},
		"utm_term":     {"page_" + slug},
	}

	o := offer{
		ID:   crs.ID,
		Name: lo.FromPtr(page.Header),
		// Encode sorts keys, matching Rails' to_query.
		URL:         y.localized("/languages/"+url.PathEscape(slug)) + "?" + utm.Encode(),
		CategoryID:  yandexCategoryID,
		CurrencyID:  currencyID,
		Params:      fixedParams,
		Description: lo.FromPtr(page.Description),
	}
	if page.LanguageCategoryID != nil {
		o.SetIDs = strconv.Itoa(*page.LanguageCategoryID)
	}

	for i, module := range plan {
		o.Plan = append(o.Plan, param{
			Name:  "План",
			Order: strconv.Itoa(i + 1),
			Unit:  lo.FromPtr(module.Name),
			Hours: planHours,
			CDATA: planIndent + lo.FromPtr(module.Description) + "\n",
		})
		for _, lv := range module.Edges.Version.Edges.LessonVersions {
			o.Plan = append(o.Plan, param{
				Name: "Ссылка на контент курса",
				Text: y.localized(strings.Join([]string{
					"/languages", url.PathEscape(lo.FromPtr(crs.Slug)),
					"lessons", url.PathEscape(lo.FromPtr(lv.Edges.Lesson.Slug)),
				}, "/")),
			})
		}
	}
	return o
}

// localized is an absolute ru page URL.
func (y *Yandex) localized(path string) string {
	return y.origin + "/" + feedLocale + path
}

type catalog struct {
	XMLName xml.Name `xml:"yml_catalog"`
	Date    string   `xml:"date,attr"`
	Shop    shop     `xml:"shop"`
}

// shop wraps sets and offers in structs rather than `sets>set` paths so the
// containers are written even when empty, as legacy wrote them.
type shop struct {
	Name        string `xml:"name"`
	Company     string `xml:"company"`
	URL         string `xml:"url"`
	Email       string `xml:"email"`
	Description string `xml:"description"`
	Picture     string `xml:"picture"`
	Sets        sets   `xml:"sets"`
	Offers      offers `xml:"offers"`
}

type sets struct {
	Set []set `xml:"set"`
}

type set struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name"`
	URL  string `xml:"url"`
}

type offers struct {
	Offer []offer `xml:"offer"`
}

// offer lists its `param` elements in two runs, split by the description.
// encoding/xml refuses two fields with the same element name, so the second
// run is an `,any` field: each entry is still written as <param> by its own
// XMLName, after the fields before it.
type offer struct {
	ID          int     `xml:"id,attr"`
	Name        string  `xml:"name"`
	URL         string  `xml:"url"`
	CategoryID  int     `xml:"categoryId"`
	Price       int     `xml:"price"`
	CurrencyID  string  `xml:"currencyId"`
	SetIDs      string  `xml:"set-ids"`
	Params      []param `xml:"param"`
	Description string  `xml:"description"`
	Plan        []param `xml:",any"`
}

// param carries either text or CDATA; an empty one of the two writes nothing.
type param struct {
	XMLName xml.Name `xml:"param"`
	Name    string   `xml:"name,attr"`
	Order   string   `xml:"order,attr,omitempty"`
	Unit    string   `xml:"unit,attr,omitempty"`
	Hours   string   `xml:"hours,attr,omitempty"`
	Text    string   `xml:",chardata"`
	CDATA   string   `xml:",cdata"`
}
