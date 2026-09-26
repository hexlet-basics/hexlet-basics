package handlers_test

import (
	"encoding/xml"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/testsupport"
)

// yandexFeed is the slice of the YML document these tests read; the full
// document is pinned by the golden test in internal/feeds.
type yandexFeed struct {
	Date   string `xml:"date,attr"`
	Offers []struct {
		Name   string `xml:"name"`
		URL    string `xml:"url"`
		Params []struct {
			Name string `xml:"name,attr"`
			Text string `xml:",chardata"`
		} `xml:"param"`
	} `xml:"shop>offers>offer"`
}

// TestYandexCoursesFeed serves the feed to a visitor at both addresses: the
// `.xml` one production answers on and the bare legacy route. The generated
// client only decodes an `application/xml` answer, so a successful call also
// pins the content type.
func TestYandexCoursesFeed(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	ctx := t.Context()

	xmlResp, err := h.Client.GetYandexCoursesFeedXml(ctx)
	require.NoError(t, err)
	fromXML, err := io.ReadAll(xmlResp)
	require.NoError(t, err)

	bareResp, err := h.Client.GetYandexCoursesFeed(ctx)
	require.NoError(t, err)
	fromBare, err := io.ReadAll(bareResp)
	require.NoError(t, err)

	for _, body := range [][]byte{fromXML, fromBare} {
		var feed yandexFeed
		require.NoError(t, xml.Unmarshal(body, &feed))
		assert.Regexp(t, `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$`, feed.Date)

		// Only JavaScript's course has three ru modules in the fixtures.
		require.Len(t, feed.Offers, 1)
		offer := feed.Offers[0]
		assert.Equal(t, "Курс JavaScript", offer.Name)
		assert.Equal(t,
			"https://code-basics.com/ru/languages/javascript-ru?utm_campaign=cb-javascript-ru"+
				"&utm_content=feed_search&utm_medium=organic&utm_source=yandex&utm_term=page_javascript-ru",
			offer.URL)

		var lessons []string
		for _, p := range offer.Params {
			if p.Name == "Ссылка на контент курса" {
				lessons = append(lessons, p.Text)
			}
		}
		assert.Equal(t, []string{
			"https://code-basics.com/ru/languages/javascript/lessons/hello-world",
			"https://code-basics.com/ru/languages/javascript/lessons/variables",
			"https://code-basics.com/ru/languages/javascript/lessons/strings",
		}, lessons)
	}
}
