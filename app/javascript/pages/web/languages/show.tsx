import { Accordion, Alert, Card, Col, Container, Row } from "react-bootstrap";

import learningEnVideo from "@/images/course-landing-page/learning_en.mp4";
import learningRuVideo from "@/images/course-landing-page/learning_ru.mp4";
import i18next from "i18next";
import { useTranslation } from "react-i18next";

import codeIllustration from "@/images/code.svg";

import XssContent from "@/components/XssContent";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { BreadcrumbItem, Language, LeadCrud, SharedProps } from "@/types";
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
import LeadFormBlock from "@/components/LeadFormBlock";

type Props = {
  lead: LeadCrud;
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
  lead,
}: Props) {
  const { t } = useTranslation();
  const {
    auth,
    locale,
  } = usePage<SharedProps>().props;

  const productSchema: WithContext<Product> = {
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
        <script type="application/ld+json">{JSON.stringify(productSchema)}</script>
      </Head>
      <Container>
        {/* {i18next.language === 'ru' && courseMember?.state === "finished" && ( */}
        {/*   <Alert variant="light" className="mb-5"> */}
        {/*     <XssContent>{t("languages.show.completed_html")}</XssContent> */}
        {/*   </Alert> */}
        {/* )} */}
        <Row className="flex-column flex-lg-row gy-5 gx-lg-5 mb-4 mb-lg-5">
          <Col className="col-lg-7">
            <div className="fs-6 fw-medium text-opacity-75 mb-3">
              {t("languages.show.free_course", { name: courseLandingPage.header })}
            </div>
            <h1 className="display-5 fw-bolder mb-3">
              {courseLandingPage.header}
            </h1>
            <div className="fs-5 text-body-secondary mb-5">
              {courseLandingPage.description}
            </div>
            <Row className="row-cols-1 gy-3">
              <Col>
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
                    className="btn btn-lg btn-primary me-3 mb-2"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      nextLesson.slug!,
                    )}
                  >
                    <span>{t("languages.show.continue")}</span>
                  </Link>
                )}
                {courseMember?.state === "finished" && (
                  <Link
                    className="btn btn-lg btn-outline-primary me-3 mb-2"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      firstLesson.slug!,
                    )}
                  >
                    <span>{t("languages.show.restart")}</span>
                  </Link>
                )}
                {courseMember?.state === "finished" && course.hexlet_program_landing_page && (
                  <a
                    target="_blank"
                    className="btn btn-lg btn-primary"
                    href={`${course.hexlet_program_landing_page}?utm_source=code-basics&utm_medium=referral`}
                  >
                    <span>{t("languages.show.hexlet_program_link")}</span>
                  </a>
                )}
              </Col>
            </Row>
          </Col>
          <Col>
            <img
              src={codeIllustration}
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
        <Row className="row-cols-1 row-cols-lg-2 gx-lg-5 gy-3 mb-lg-4 py-3 py-md-5 justify-content-between">
          <Col className="order-last order-lg-first">
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
        <div className="mb-lg-5 py-5">
          <div className="display-5 fw-semibold lh-1 mb-5">
            {t("languages.show.learning_program")}
          </div>
          <Accordion
            defaultActiveKey={["0"]}
            className="hexlet-basics-accordion"
          >
            <Row className="gx-lg-5">
              <Col>
                {courseModules.map((m, index) => (
                  <Accordion.Item
                    eventKey={index.toString()}
                    className="rounded-0 border-0 border-bottom border-secondary-subtle py-3"
                    key={m.id}
                  >
                    <Accordion.Header as="h3">{m.name!}</Accordion.Header>
                    <Accordion.Body className="px-0 pb-0">
                      <Row>
                        <Col className="col-4">
                          <ul className="">
                            {(lessonsByModuleId[m.id] ?? []).map((l) => (
                              <li key={l.id}>
                                <Link
                                  className="text-decoration-none link-body-emphasis"
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
                        </Col>
                        <Col>
                          {m.description}
                        </Col>
                      </Row>
                    </Accordion.Body>
                  </Accordion.Item>
                ))}
              </Col>
            </Row>
          </Accordion>
        </div>
        <div className="pb-5 pt-3 mb-lg-5">
          <Card className="bg-primary border-0 py-5">
            <Row className="justify-content-center my-5">
              <Col className="col-11 col-lg-7">
                <div className="text-light display-5 fw-semibold lh-1 text-center mb-5">
                  {t("languages.show.course_graduates")}
                </div>
                <div className="d-flex flex-column flex-sm-row justify-content-center gap-2">
                  <Link
                    className="btn btn-light btn-lg text-primary"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      firstLesson.slug,
                    )}
                  >
                    <span>{t("languages.show.start")}</span>
                  </Link>
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
            <div className="pe-lg-5">
              <div className="d-flex mb-3">
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
                <div>
                  <div className="fw-bold">
                    {t("languages.show.convenient format")}
                  </div>
                  <XssContent>
                    {t("languages.show.learning_conveniently")}
                  </XssContent>
                </div>
              </div>
              <div className="d-flex mb-3">
                <div>
                  <div className="fw-bold">
                    {t("languages.show.browser_practice")}
                  </div>
                  <XssContent>
                    {t("languages.show.real_life_challenges")}
                  </XssContent>
                </div>
              </div>
              <div className="d-flex">
                <div>
                  <div className="fw-bold">
                    {t("languages.show.ai_without_limits")}
                  </div>
                  <XssContent>{t("languages.show.ai_explanation")}</XssContent>
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
        <div className="py-4 mb-lg-5">
          <Card className="px-3 px-sm-4 px-lg-5 py-5">
            <Row className="gy-4 justify-content-between align-items-center py-lg-5">
              <Col className="col-12 col-lg-9 col-xxl-8">
                <div className="display-6 fw-semibold lh-1">
                  {t("languages.show.ai_learning")}
                </div>
              </Col>
              <Col className="col-auto">
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
            </Row>
          </Card>
        </div>
        {/* <div className="py-4 py-lg-5">
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
        </div> */}
        {i18next.language === "ru" && (
          <div className="py-4 py-lg-5">
            <div className="bg-dark text-light rounded-3 overflow-hidden">
              <Row className="flex-column flex-lg-row justify-content-end m-0">
                <Col className="position-relative p-0">
                  <div className="hexlet-basics-community-image" />
                </Col>
                <Col className="col-lg-7 col-xl-6 p-4 p-md-5">
                  <div className="d-flex flex-column justify-content-center py-3 py-lg-4">
                    <div className="display-5 fw-semibold lh-1 mb-4">
                      {t("languages.show.more_than_support")}
                    </div>
                    <div className="pe-lg-5">
                      <div className="mb-5 pe-xl-4">
                        {t("languages.show.about_developer_community")}
                      </div>
                      <a
                        className="btn btn-secondary"
                        href="https://t.me/HexletLearningBot"
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        <span>{t("languages.show.join")}</span>
                      </a>
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
