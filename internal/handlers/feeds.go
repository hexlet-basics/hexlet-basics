package handlers

import (
	"bytes"
	"context"

	"hexletbasics/internal/api"
)

// GetYandexCoursesFeedXml serves the Yandex course catalogue at
// `/api/feeds/yandex_courses.xml`, the address production answers on and
// Yandex is registered with.
func (s *Server) GetYandexCoursesFeedXml(ctx context.Context) (api.GetYandexCoursesFeedXmlOK, error) {
	doc, err := s.yandexFeed.Build(ctx)
	if err != nil {
		return api.GetYandexCoursesFeedXmlOK{}, err
	}
	return api.GetYandexCoursesFeedXmlOK{Data: bytes.NewReader(doc)}, nil
}

// GetYandexCoursesFeed serves the same catalogue at the bare path the legacy
// route declares (which legacy itself answered with 406).
func (s *Server) GetYandexCoursesFeed(ctx context.Context) (api.GetYandexCoursesFeedOK, error) {
	doc, err := s.yandexFeed.Build(ctx)
	if err != nil {
		return api.GetYandexCoursesFeedOK{}, err
	}
	return api.GetYandexCoursesFeedOK{Data: bytes.NewReader(doc)}, nil
}
