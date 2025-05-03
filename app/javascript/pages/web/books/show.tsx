import type { PropsWithChildren } from "react";
import * as Routes from "@/routes.js";
import { Col, Container, Row } from "react-bootstrap";

import bookCoverImg from "@/images/profession-developer-book-cover.png";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout.tsx";
import { useTranslation } from "react-i18next";
import { Link } from "@inertiajs/react";
import bookToc from '@/lib/book.ts';

type Props = PropsWithChildren & {
  bookRequested: boolean
  // user: User;
};

export default function New({ bookRequested }: Props) {
  const { t } = useTranslation()
  return (
    <ApplicationLayout>
      <Container>
        <div className="p-5 mt-3">
          <Row>
            <Col className="col-12 col-md-7">
              <div className="text-muted fs-5">{t('books.show.freebook')}</div>
              <h1>{t('books.show.header')}</h1>
              <p className="lead">{t('books.show.description')}</p>
              {!bookRequested && <Link
                method="post"
                href={Routes.create_request_book_url()}
                className="btn btn-lg btn-outline-primary"
              >
                {t('books.show.request')}
              </Link>}
              {bookRequested && <a
                target="_blank"
                href={Routes.download_book_url()}
                className="btn btn-lg btn-outline-primary"
              >
                {t('books.show.download')}
              </a>}
              <Row className="mt-5">
                <Col>
                  <div className="d-flex">
                    <i className="bi bi-compass me-2 my-auto" />
                    <span className="fw-bold fs-5">{t('books.show.features.direction')}</span>
                  </div>
                  <div>{t('books.show.features.direction_explanation')} </div>
                </Col>
                <Col>
                  <div className="d-flex">
                    <i className="bi bi-compass me-2 my-auto" />
                    <span className="fw-bold fs-5">{t('books.show.features.plan')}</span>
                  </div>
                  <div>{t('books.show.features.plan_explanation')} </div>
                </Col>
              </Row>
              <Row className="mt-4">
                <Col>
                  <div className="d-flex">
                    <i className="bi bi-compass me-2 my-auto" />
                    <span className="fw-bold fs-5">{t('books.show.features.resume')}</span>
                  </div>
                  <div>{t('books.show.features.resume_explanation')} </div>
                </Col>
                <Col>
                  <div className="d-flex">
                    <i className="bi bi-compass me-2 my-auto" />
                    <span className="fw-bold fs-5">{t('books.show.features.interview')}</span>
                  </div>
                  <div>{t('books.show.features.interview_explanation')} </div>
                </Col>
              </Row>
            </Col>
            <Col>
              <img src={bookCoverImg} className="img-fluid" />
            </Col>
          </Row>
          <h2 className="mb-4">{t('books.show.toc')}</h2>
          <Row className="flex-column">
            {bookToc.map((item, index) => (
              <Col key={item.title} className="py-3">
                <Row>
                  <Col className="col-12 col-md-2 mb-3 fw-bold">
                    {t('books.show.chapter', { number: index + 1 })}
                  </Col>
                  <Col className="col-12 col-md-4 mb-3">{item.title}</Col>
                  <Col className="col-12 col-md-6 mb-3">
                    <ul>
                      {item.subsections.map((subsection, index) => (
                        <li key={index}>{subsection}</li>
                      ))}
                    </ul>
                  </Col>
                </Row>
              </Col>
            ))}
          </Row>
        </div>
      </Container>
    </ApplicationLayout>
  );
}

