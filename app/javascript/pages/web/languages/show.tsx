import { Accordion, Alert, Card, Col, Container, Row } from "react-bootstrap";

import learningEnVideo from "@/images/course-landing-page/learning_en.mp4";
import learningRuVideo from "@/images/course-landing-page/learning_ru.mp4";
import jsImage from "@/images/javascript.png";
import i18next from "i18next";
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
  Review,
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
  reviews: Review[];
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
  reviews,
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
        <Row className="flex-column flex-lg-row gy-4 gy-md-5 gx-lg-5 mb-4 mb-lg-5">
          <Col className="col-lg-7">
            <div className="fs-5 fw-medium text-primary text-opacity-75 mb-3">
              {t("languages.show.free_course")}
            </div>
            <h1 className="display-5 fw-bolder mb-3">
              {courseLandingPage.header}
            </h1>
            <div className="fs-5 text-body-secondary mb-5">
              {courseLandingPage.description}
            </div>
            <Row className="row-cols-1 row-cols-sm-auto gy-3">
              {!courseMember && (
                <Col className="col-lg-5 col-xl-4">
                  <Link
                    className="btn btn-lg btn-primary"
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
                <Col>
                  <Link
                    className="btn btn-lg btn-outline-primary"
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
                <Col>
                  <Link
                    className="btn btn-lg fw-medium"
                    href={Routes.new_user_path()}
                  >
                    <span className="me-2">
                      {t("languages.show.registration")}
                    </span>
                    <i className="bi bi-arrow-right" />
                  </Link>
                </Col>
              )}
            </Row>
          </Col>
          <Col>
            <img
              // src={course.cover_list_variant}
              src={jsImage}
              width="100%"
              height="auto"
              alt={t("languages.show.cover_image")}
            />
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
        <div className="pb-5 mb-lg-5">
          <Card className="bg-primary border-0 py-5">
            <Row className="justify-content-center my-5">
              <Col className="col-11 col-lg-7">
                <div className="text-light display-5 fw-semibold lh-1 text-center mb-5">
                  {t("languages.show.course_graduates")}
                </div>
                <div className="d-flex flex-column flex-sm-row justify-content-center gap-2">
                  {!courseMember && (
                    <Link
                      className="btn btn-light btn-lg text-primary"
                      href={Routes.language_lesson_path(
                        course.slug!,
                        firstLesson.slug,
                      )}
                    >
                      <span>{t("languages.show.start")}</span>
                    </Link>
                  )}
                  {courseMember && nextLesson && (
                    <Link
                      className="btn btn-lg btn-outline-light"
                      href={Routes.language_lesson_path(
                        course.slug!,
                        nextLesson.slug,
                      )}
                    >
                      <span>{t("languages.show.continue")}</span>
                    </Link>
                  )}
                  {auth.user.guest && (
                    <Link
                      className="btn btn-lg text-light fw-medium"
                      href={Routes.new_user_path()}
                    >
                      <span className="me-2">
                        {t("languages.show.registration")}
                      </span>
                      <i className="bi bi-arrow-right" />
                    </Link>
                  )}
                </div>
              </Col>
            </Row>
          </Card>
        </div>
        <div className="display-5 fw-semibold lh-1 mb-4 mb-lg-0">
          {t("languages.show.about_learning")}
        </div>
        <Row className="row-cols-1 row-cols-lg-2 mb-lg-5 pb-4 py-md-5 gy-4">
          <Col>
            <div className="me-lg-5 pe-lg-4">
              <div className="d-flex mb-3">
                <i className="bi bi-cloud-arrow-up-fill me-3 text-primary" />
                <div>
                  <div className="fw-bold">
                    {t("languages.show.try_without_registration")}
                  </div>
                  <XssContent>
                    {t("languages.show.without_registration")}
                  </XssContent>
                </div>
              </div>
              <div className="d-flex mb-3">
                <i className="bi bi-laptop-fill me-3 text-primary" />
                <div>
                  <div className="fw-bold">
                    {t("languages.show.convenient format")}
                  </div>
                  <XssContent>
                    {t("languages.show.learning_conveniently")}
                  </XssContent>
                </div>
              </div>
              <div className="d-flex">
                <i className="bi bi-lock-fill me-3 text-primary" />
                <div>
                  <div className="fw-bold">
                    {t("languages.show.browser_practice")}
                  </div>
                  <XssContent>
                    {t("languages.show.real_life_challenges")}
                  </XssContent>
                </div>
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
        <Row className="py-4 mb-lg-5">
          <Col className="col-lg-9">
            <div className="display-5 fw-semibold lh-1 mb-4">
              {t("languages.show.about_ai")}
            </div>
            <div className="text-body-secondary">
              <div className="mb-3">{t("languages.show.support_from_ai")}</div>
              <div className="fw-bold mb-2">
                {t("languages.show.ai_features")}
              </div>
              <ul>
                <li>{t("languages.show.explain_topics")}</li>
                <li>{t("languages.show.help_to_understand")}</li>
                <li>{t("languages.show.answers_questions")}</li>
                <li>{t("languages.show.always_in_touch")}</li>
              </ul>
            </div>
          </Col>
        </Row>
        <div className="py-4 mb-lg-5">
          <Card className="px-3 px-sm-4 px-lg-5 py-5">
            <Row className="gy-4 justify-content-between align-items-center py-lg-5">
              <Col className="col-12 col-lg-9 col-xxl-8">
                <div className="display-6 fw-semibold lh-1">
                  {t("languages.show.ai_learning")}
                </div>
              </Col>
              <Col className="col-auto">
                {!courseMember && (
                  <Link
                    className="btn btn-lg btn-primary"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      firstLesson.slug,
                    )}
                  >
                    <span>{t("languages.show.start")}</span>
                  </Link>
                )}
                {courseMember && nextLesson && (
                  <Link
                    className="btn btn-lg btn-primary"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      nextLesson.slug,
                    )}
                  >
                    <span>{t("languages.show.continue")}</span>
                  </Link>
                )}
              </Col>
            </Row>
          </Card>
        </div>
        <div className="py-4 py-lg-5">
          <div className="display-5 fw-semibold lh-1 mb-5">
            {t("languages.show.reviews")}
          </div>
          <Row className="row-cols-1 row-cols-lg-2 row-cols-xl-3 gy-4">
            {reviews.map((review) => (
              <Col key={review.id}>
                <div className="d-flex flex-column h-100 rounded-3 bg-body-tertiary p-4">
                  <div className="mb-3">{review.body}</div>
                  <div className="mt-auto">
                    <div className="fw-bold">{review.user.name}</div>
                    {review.language.slug && (
                      <Link
                        className="text-decoration-none text-body"
                        href={Routes.language_path(courseLandingPage.slug)}
                      >
                        <span>{review.language.slug}</span>
                      </Link>
                    )}
                  </div>
                </div>
              </Col>
            ))}
          </Row>
        </div>
        {i18next.language === "ru" && (
          <div className="py-4 py-lg-5">
            <div className="bg-dark text-light rounded-3 overflow-hidden">
              <Row className="flex-column flex-lg-row justify-content-end m-0">
                <Col className="position-relative p-0">
                  <div className="hexlet-basics-community-image" />
                </Col>
                <Col className="col-lg-7 col-xl-6 p-4 p-lg-5">
                  <div className="d-flex flex-column justify-content-center py-3 py-lg-4">
                    <div className="display-5 fw-semibold lh-1 mb-4">
                      {t("languages.show.more_than_support")}
                    </div>
                    <div className="pe-lg-5">
                      <div className="mb-5">
                        {t("languages.show.about_developer_community")}
                      </div>
                      <Link
                        className="btn btn-secondary"
                        href="https://ttttt.me/HexletLearningBot"
                      >
                        <span>{t("languages.show.join")}</span>
                      </Link>
                    </div>
                  </div>
                </Col>
              </Row>
            </div>
          </div>
        )}
        {courseLandingPageQnaItems.length > 0 && (
          <div className="py-4 py-lg-5">
            <div className="display-5 fw-semibold lh-1 mb-5">
              {t("languages.show.sort_questions")}
            </div>
            <Row className="gy-4 gy-md-5">
              {courseLandingPageQnaItems.map((item) => (
                <Col key={item.id} className="col-12 col-md-6">
                  <div className="fs-5 fw-medium mb-3 pe-lg-5">
                    {item.question}
                  </div>
                  <p className="pe-lg-5">{item.answer}</p>
                </Col>
              ))}
            </Row>
          </div>
        )}
      </Container>
    </ApplicationLayout>
  );
}
