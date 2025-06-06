import type { PropsWithChildren } from "react";
import * as Routes from "@/routes.js";
import { Col, Container, Row } from "react-bootstrap";

import bookCoverImg from "@/images/profession-developer-book-cover.png";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout.tsx";
import { useTranslation } from "react-i18next";
import { Link, usePage } from "@inertiajs/react";
import bookToc from '@/lib/book.ts';
import i18next from "i18next";
import { LeadCrud, SharedProps } from "@/types";
import LeadFormBlock from "@/components/LeadFormBlock";

type Props = PropsWithChildren & {
  bookRequested: boolean
  lead: LeadCrud
  // user: User;
};

export default function Show({ bookRequested, lead }: Props) {
  const { t } = useTranslation()
  const {
    auth,
  } = usePage<SharedProps>().props;

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
              <Row className="mt-5 row-cols-1 row-cols-md-2">
                <Col className="mb-4">
                  <div className="d-flex">
                    <i className="bi bi-compass me-2 my-auto" />
                    <span className="fw-bold fs-5">{t('books.show.features.direction')}</span>
                  </div>
                  <div>{t('books.show.features.direction_explanation')} </div>
                </Col>
                <Col className="mb-4">
                  <div className="d-flex">
                    <i className="bi bi-compass me-2 my-auto" />
                    <span className="fw-bold fs-5">{t('books.show.features.plan')}</span>
                  </div>
                  <div>{t('books.show.features.plan_explanation')} </div>
                </Col>
              </Row>
              <Row className="row-cols-1 row-cols-md-2">
                <Col className="mb-4">
                  <div className="d-flex">
                    <i className="bi bi-compass me-2 my-auto" />
                    <span className="fw-bold fs-5">{t('books.show.features.resume')}</span>
                  </div>
                  <div>{t('books.show.features.resume_explanation')} </div>
                </Col>
                <Col className="mb-4">
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
          {bookToc.map((item, index) => (
            <div key={item.title}>
              <Row className="py-3 bottom">
                <Col className="col-12 col-md-2 fw-bold">
                  {t('books.show.chapter', { number: index + 1 })}
                </Col>
                <Col className="col-12 col-md-4">{item.title}</Col>
                <Col className="col-12 col-md-6">
                  <ul>
                    {item.subsections.map((subsection, index) => (
                      <li key={index}>{subsection}</li>
                    ))}
                  </ul>
                </Col>
              </Row>
              {(index !== bookToc.length - 1) && <hr />}
            </div>
          ))}

        </div>

        {!auth.user.guest && i18next.language === 'ru' && (
          <Row className="align-items-center py-5">
            <Col className="col-12 col-lg-7 fw-bold display-4 mb-5">
              {t("home.index.consultation")}
            </Col>
            <Col className="col-12 col-lg-5 mx-auto d-flex flex-column">
              <div className="bg-body-tertiary p-4 p-md-5 border rounded-3">
                <LeadFormBlock lead={lead} />
              </div>
            </Col>
          </Row>
        )}
      </Container>
    </ApplicationLayout>
  );
}
