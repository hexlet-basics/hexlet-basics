import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type {
  LanguageCategory,
  LanguageLandingPage,
  LanguageLandingPageForLists, LanguageLandingPageQnaItem, LeadCrud
} from "@/types";

import CourseBlock from "@/components/CourseBlock";
import * as Routes from "@/routes.js";
import i18next from "i18next";
import { usePage } from "@inertiajs/react";
import { SharedProps } from "@/types";
import LeadFormBlock from "@/components/LeadFormBlock";
import MarkdownViewer from "@/components/MarkdownViewer";

type Props = PropsWithChildren & {
  categoryLandingPages: LanguageLandingPageForLists[];
  qnaItems: LanguageLandingPageQnaItem[];
  landingPages: LanguageLandingPage[];
  courseCategory: LanguageCategory;
  lead: LeadCrud;
};

export default function Show({
  courseCategory,
  qnaItems,
  categoryLandingPages,
  lead,
}: Props) {
  const { t } = useTranslation();

  const {
    auth: { user },
  } = usePage<SharedProps>().props;

  const items = [
    {
      name: t("language_categories.index.header"),
      url: Routes.language_categories_path(),
    },
    {
      name: courseCategory.header!,
      url: Routes.language_category_path(courseCategory.slug!),
    },
  ];

  return (
    <ApplicationLayout items={items} header={courseCategory.header!}>
      <Container>

        {courseCategory.description && (
          <Row className="mb-5">
            <Col className="col-12 col-sm-8">
              <p className="fs-5">{courseCategory.description}</p>
            </Col>
          </Row>
        )}

        <Row className="row-cols-2 row-cols-md-3 row-cols-lg-4">
          {categoryLandingPages.map((lp) => (
            <Col className="mb-4" key={lp.id}>
              <CourseBlock landingPage={lp} />
            </Col>
          ))}
        </Row>

        {!user.guest && i18next.language === 'ru' && (
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

        {qnaItems.length > 0 && (
          <div className="py-4 py-lg-5">
            <div className="display-5 fw-semibold lh-1 mb-5">
              {t("languages.show.sort_questions")}
            </div>
            <Row className="gy-4 gy-md-5">
              {qnaItems.map((item) => (
                <Col key={item.id} className="col-12 col-md-6">
                  <div className="fs-5 fw-medium mb-3 pe-lg-5">
                    {item.question}
                  </div>
                  <p className="pe-lg-5"><MarkdownViewer>{item.answer}</MarkdownViewer></p>
                </Col>
              ))}
            </Row>
          </div>
        )}

      </Container>

    </ApplicationLayout>
  );
}
