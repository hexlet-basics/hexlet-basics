import { Accordion, Alert, Col, Container, Row } from "react-bootstrap";

import learningEnVideo from "@/images/course-landing-page/learning_en.mp4";
import learningRuVideo from "@/images/course-landing-page/learning_ru.mp4";
import { useTranslation } from "react-i18next";

import XssContent from "@/components/XssContent";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { BreadcrumbItem, Language, SharedProps } from "@/types";
import type {
  LanguageCategory,
  LanguageLandingPage,
  LanguageLandingPageQnaItem,
  LanguageLesson,
  LanguageMember,
  LanguageModule,
} from "@/types/serializers";
import { Head, Link, usePage } from "@inertiajs/react";
import type { Product, WithContext } from "schema-dts";

type Props = {
  course: Language;
  courseMember?: LanguageMember;
  courseCategory: LanguageCategory;
  courseLandingPage: LanguageLandingPage;
  courseLandingPageQnaItems: LanguageLandingPageQnaItem[];
  firstLesson: LanguageLesson;
  nextLesson?: LanguageLesson;
  // user: User;
  courseModules: LanguageModule[];
  lessonsByModuleId: {
    [moduleId: number]: LanguageLesson[];
  };
};

export default function Show({
  firstLesson,
  nextLesson,
  courseLandingPage,
  courseLandingPageQnaItems,
  courseMember,
  courseCategory,
  courseModules,
  lessonsByModuleId,
  course,
}: Props) {
  const { t } = useTranslation();
  const { auth, locale } = usePage<SharedProps>().props;

  const product: WithContext<Product> = {
    "@context": "https://schema.org",
    "@type": "Product",
    description: courseLandingPage.description,
    image: course.cover_list_variant,
    offers: {
      "@type": "Offer",
      price: "0",
      priceCurrency: "RUB",
      availability: "https://schema.org/InStock",
    },
    aggregateRating: {
      "@type": "AggregateRating",
      ratingValue: course.rating_value,
      ratingCount: course.rating_count,
    },
    name: courseLandingPage.header,
    // TODO: add review
  };

  const breadcrumbItems: BreadcrumbItem[] = [
    {
      name: courseCategory.name,
      url: Routes.language_category_url(courseCategory.slug!),
    },
    {
      name: courseLandingPage.header,
      url: Routes.language_url(courseLandingPage.slug),
    },
  ];

  return (
    <ApplicationLayout items={breadcrumbItems}>
      <Head>
        <script type="application/ld+json">{JSON.stringify(product)}</script>
      </Head>
      <Container className="pt-4 pt-lg-5">
        {courseMember?.state === "finished" && (
          <Alert variant="success">
            <XssContent>{t("languages.show.completed_html")}</XssContent>
          </Alert>
        )}
        <Row className="justify-content-center py-2 mb-5">
          <Col className="col-lg-7 text-center">
            <div className="fs-5 fw-medium text-primary text-opacity-75 lh-sm mb-2">
              {t("languages.show.free_course")}
            </div>
            <h1 className="display-5 fw-bolder mb-3">
              {courseLandingPage.header}
            </h1>
            <div className="fs-5 text-body-secondary mb-5">
              {courseLandingPage.description}
            </div>
            <Row className="row-cols-1 row-cols-sm-2 justify-content-center gy-3">
              {!courseMember && (
                <Col className="col-lg-4">
                  <Link
                    className="btn btn-primary w-100"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      firstLesson.slug,
                    )}
                  >
                    <span>{t("languages.show.start")}</span>
                  </Link>
                </Col>
              )}
              {courseMember && nextLesson && (
                <Col className="col-lg-4">
                  <Link
                    className="btn btn-outline-primary w-100"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      nextLesson.slug,
                    )}
                  >
                    <span>{t("languages.show.continue")}</span>
                  </Link>
                </Col>
              )}
              {auth.user.guest && (
                <Col className="col-lg-5">
                  <Link className="btn fw-medium" href={Routes.new_user_path()}>
                    <span className="me-2">
                      {t("languages.show.registration")}
                    </span>
                    <i className="bi bi-arrow-right" />
                  </Link>
                </Col>
              )}
            </Row>
          </Col>
        </Row>
        <Row className="mb-lg-5 py-4 py-md-5">
          <Col className="col-lg-9">
            <h2 className="display-5 fw-medium lh-1 mb-4">
              {courseLandingPage.used_in_header}
            </h2>
            <p className="text-body-secondary pe-lg-2">
              {courseLandingPage.used_in_description}
            </p>
          </Col>
        </Row>
        <Row className="row-cols-1 row-cols-md-2 gx-lg-5 gy-5 mb-lg-4 py-3 py-md-5 justify-content-between">
          <Col>
            {courseLandingPage.outcomes_image && (
              <img
                src={courseLandingPage.outcomes_image}
                width="100%"
                height="auto"
                alt={t("languages.show.learning_preview")}
                className="rounded-5 shadow-lg"
              />
            )}
          </Col>
          <Col>
            <div className="display-5 fw-semibold lh-1 mb-4">
              {courseLandingPage.outcomes_header}
            </div>
            <p className="text-body-secondary pe-lg-5">
              {courseLandingPage.outcomes_description}
            </p>
          </Col>
        </Row>
        <Row className="mb-lg-5 py-5">
          <Col className="col-lg-10">
            <div className="display-5 fw-semibold lh-1 mb-4">
              {t("languages.show.learning_program")}
            </div>
            <Accordion defaultActiveKey="0" className="hexlet-basics-accordion">
              {courseModules.map((m, index) => (
                <Accordion.Item
                  eventKey={index.toString()}
                  key={m.id}
                  className="rounded-0 border-0 border-bottom border-secondary-subtle py-3 py-md-4"
                >
                  <Accordion.Header as="h3">
                    {index + 1}. {m.name!}
                  </Accordion.Header>
                  <Accordion.Body className="px-0 pb-0">
                    <ul className="list-unstyled">
                      {(lessonsByModuleId[m.id] ?? []).map((l) => (
                        <li key={l.id}>
                          <Link
                            className="text-decoration-none text-body-secondary"
                            href={Routes.language_lesson_path(
                              course.slug!,
                              l.slug!,
                            )}
                          >
                            <span>{l.name}</span>
                          </Link>
                        </li>
                      ))}
                    </ul>
                  </Accordion.Body>
                </Accordion.Item>
              ))}
            </Accordion>
          </Col>
        </Row>
        <div className="display-5 fw-semibold lh-1 mb-4 mb-lg-0">
          {t("languages.show.about_learning")}
        </div>
        <Row className="row-cols-1 row-cols-lg-2 mb-lg-5 pb-4 py-md-5 gy-4">
          <Col>
            <div className="me-lg-5 pe-lg-4">
              <div className="d-flex mb-3">
                <i className="bi bi-cloud-arrow-up-fill me-3 text-primary" />
                <XssContent>
                  {t("languages.show.without_registration_html")}
                </XssContent>
              </div>
              <div className="d-flex mb-3">
                <i className="bi bi-laptop-fill me-3 text-primary" />
                <XssContent>
                  {t("languages.show.learning_conveniently_html")}
                </XssContent>
              </div>
              <div className="d-flex">
                <i className="bi bi-lock-fill me-3 text-primary" />
                <XssContent>
                  {t("languages.show.real_life_challenges_html")}
                </XssContent>
              </div>
            </div>
          </Col>
          <Col>
            <div className="bg-primary rounded-5 overflow-hidden py-lg-3">
              <video
                className="w-100 rounded-4 hexlet-basics-learning-video"
                src={locale === "en" ? learningEnVideo : learningRuVideo}
                autoPlay
                loop
                muted
                playsInline
              />
            </div>
          </Col>
        </Row>
        {courseLandingPageQnaItems.length > 0 && (
          <Row className="mb-lg-5 py-4 py-lg-5">
            <Col className="col-lg-10">
              <div className="display-5 fw-semibold lh-1 mb-4">
                {t("languages.show.sort_questions")}
              </div>
              <Accordion
                defaultActiveKey="0"
                className="hexlet-basics-accordion"
              >
                {courseLandingPageQnaItems.map((item, index) => (
                  <Accordion.Item
                    eventKey={index.toString()}
                    key={item.id}
                    className="rounded-0 border-0 border-bottom border-secondary-subtle py-3 py-md-4"
                  >
                    <Accordion.Header as="h3">{item.question}</Accordion.Header>
                    <Accordion.Body className="px-0">
                      {item.answer}
                    </Accordion.Body>
                  </Accordion.Item>
                ))}
              </Accordion>
            </Col>
          </Row>
        )}
      </Container>
    </ApplicationLayout>
  );
}
