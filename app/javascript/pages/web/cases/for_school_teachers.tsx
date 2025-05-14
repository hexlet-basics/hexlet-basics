import learningEnVideo from "@/images/course-landing-page/learning_en.mp4";
import learningRuVideo from "@/images/course-landing-page/learning_ru.mp4";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";
import { Link } from "@inertiajs/react";
import i18next from "i18next";
import { Col, Container, Row } from "react-bootstrap";
import { useTranslation } from "react-i18next";

export default function ForSchoolTeachersBlock() {
  const {
    auth: { user },
  } = usePage<SharedProps>().props;
  const { t } = useTranslation();

  const programmingBasicList = t("for_school_teachers.programming_basic_list", {
    returnObjects: true,
  });

  const interactiveApproachList = t(
    "for_school_teachers.interactive_approach_list",
    { returnObjects: true },
  );

  const earlyCareerGuidanceList = t(
    "for_school_teachers.early_career_guidance_list",
    { returnObjects: true },
  );

  const howToLearnCards = t(
    "for_school_teachers.how_to_learn_programming_cards",
    { returnObjects: true },
  );

  return (
    <ApplicationLayout>
      <Container className="py-3">
        <Row className="row-cols-1 row-cols-lg-2 gy-5 gy-lg-0 align-items-center py-3 py-lg-5">
          <Col>
            <h1 className="display-5 fw-normal lh-1 mb-4 mb-lg-5">
              {t("for_school_teachers.header")}
            </h1>
            <div className="fs-5 lh-1 text-body-secondary mb-4 mb-lg-5 pb-lg-3">
              {t("for_school_teachers.description")}
            </div>
            <a
              className="btn btn-lg btn-primary w-100 w-lg-auto px-lg-5"
              href={user.guest ? Routes.new_user_path() : "https://ru.hexlet.io/courses"}
            >
              <span>{t("for_school_teachers.try")}</span>
            </a>
          </Col>
          <Col>
            <div className="hexlet-basics-cases-block-image" />
          </Col>
        </Row>
        <Row className="row-cols-1 row-cols-lg-2 py-5 gy-3">
          <Col>
            <div className="h1 fw-normal lh-1">
              {t("for_school_teachers.integrate_into_education")}
            </div>
          </Col>
          <Col>
            <p className="fs-5 text-body-secondary lh-1 pe-lg-4">
              {t("for_school_teachers.lay_programming_foundations")}
            </p>
          </Col>
        </Row>
        <Row className="row-cols-1 row-cols-lg-2 py-2 gy-3">
          <Col>
            <div className="pe-lg-5">
              <div className="fs-5 fw-medium lh-sm">
                {t("for_school_teachers.programming_competently")}
              </div>
              <ul className="text-body-secondary lh-1 list-unstyled my-4 pb-2 py-lg-2">
                {Array.isArray(programmingBasicList) &&
                  programmingBasicList.map((item, index) => (
                    <li key={index.toString()} className="d-flex gap-3 mb-3">
                      <span>–</span>
                      <span>{t(item, { defaultValue: item })}</span>
                    </li>
                  ))}
              </ul>
            </div>
            <div className="pe-lg-5">
              <div className="fs-5 fw-medium lh-sm">
                {t("for_school_teachers.interactive_approach")}
              </div>
              <ul className="text-body-secondary lh-1 list-unstyled my-4 pb-2 py-lg-2">
                {Array.isArray(interactiveApproachList) &&
                  interactiveApproachList.map((item, index) => (
                    <li key={index.toString()} className="d-flex gap-3 mb-3">
                      <span>–</span>
                      <span>{t(item, { defaultValue: item })}</span>
                    </li>
                  ))}
              </ul>
            </div>
            <div className="pe-lg-5">
              <div className="fs-5 fw-medium lh-sm">
                {t("for_school_teachers.early_career_guidance")}
              </div>
              <ul className="text-body-secondary lh-1 list-unstyled my-4 py-lg-2">
                {Array.isArray(earlyCareerGuidanceList) &&
                  earlyCareerGuidanceList.map((item, index) => (
                    <li key={index.toString()} className="d-flex gap-3 mb-3">
                      <span>–</span>
                      <span>{t(item, { defaultValue: item })}</span>
                    </li>
                  ))}
              </ul>
            </div>
          </Col>
          <Col>
            <div className="bg-primary rounded-5 overflow-hidden py-lg-3">
              <video
                className="w-100 rounded-4 hexlet-basics-learning-video"
                src={
                  i18next.language === "en" ? learningEnVideo : learningRuVideo
                }
                autoPlay
                loop
                muted
                playsInline
              />
            </div>
          </Col>
        </Row>
        <div className="py-4 py-lg-5 mt-4 mt-lg-0">
          <div className="bg-dark text-light rounded-4 overflow-hidden p-3 p-lg-4">
            <div className="h1 lh-1 fw-normal w-lg-75 pb-3 pe-lg-5">
              {t("for_school_teachers.integrate_now")}
            </div>
            <div className="d-flex flex-column flex-lg-row justify-content-between align-items-lg-end gap-5 gap-lg-0">
              <div className="text-white-50 lh-1">
                {t("for_school_teachers.open_browser_and_sign_up")}
              </div>
              {user.guest && (
                <Link
                  className="btn btn-lg btn-light w-100 w-lg-auto px-lg-5"
                  href={Routes.new_user_path()}
                >
                  <span>{t("for_school_teachers.sign_up")}</span>
                </Link>
              )}
              {!user.guest && (
                <a
                  className="btn btn-lg btn-light w-100 w-lg-auto px-lg-5"
                  href="https://ru.hexlet.io/courses"
                  target="_blank"
                  rel="noreferrer"
                >
                  <span>{t("for_school_teachers.select_course")}</span>
                </a>
              )}
            </div>
          </div>
        </div>
        <div className="py-3 py-lg-5">
          <div className="h1 fw-normal lh-1 mb-4 mb-md-5 w-lg-75 w-xxl-50 pe-lg-4">
            {t("for_school_teachers.how_to_learn_programming")}
          </div>
          <Row className="row-cols-1 row-cols-md-2 row-cols-xl-4 gy-4">
            {Array.isArray(howToLearnCards) &&
              howToLearnCards.map((item, index) => (
                <Col key={index.toString()}>
                  <div className="border rounded-4 p-3 p-lg-4 h-100 d-flex flex-column gap-4">
                    <img
                      src={`https://code-basics.test/vite-dev/images/for-school-teachers-page/${t(item.img, { defaultValue: item.img })}.svg`}
                      width="90"
                      height="65"
                      alt={t(item.img, { defaultValue: item.img })}
                    />
                    <div className="lh-sm">
                      {t(item.title, { defaultValue: item.title })}
                    </div>
                    <div className="lh-sm">
                      {t(item.subtitle, { defaultValue: item.subtitle })}
                    </div>
                  </div>
                </Col>
              ))}
          </Row>
        </div>
        <div className="pt-5 pb-3 pb-lg-5">
          <Row className="bg-dark text-light rounded-4 d-flex align-items-end justify-content-between p-2 p-md-3 p-lg-4 m-0 gy-2">
            <Col className="col-md-10 col-lg-7">
              <div className="h1 lh-1 fw-normal pe-lg-4 mb-4">
                {t("for_school_teachers.sign_up_and_start_learning")}
              </div>
            </Col>
            <Col className="col-lg-auto">
              {user.guest && (
                <Link
                  className="btn btn-lg btn-light w-100 w-lg-auto px-lg-5 mb-3 mb-lg-0"
                  href={Routes.new_user_path()}
                >
                  <span>{t("for_school_teachers.sign_up")}</span>
                </Link>
              )}
              {!user.guest && (
                <a
                  className="btn btn-lg btn-light w-100 w-lg-auto px-lg-5 mb-3 mb-lg-0"
                  href="https://ru.hexlet.io/courses"
                  target="_blank"
                  rel="noreferrer"
                >
                  <span>{t("for_school_teachers.select_course")}</span>
                </a>
              )}
            </Col>
          </Row>
        </div>
      </Container>
    </ApplicationLayout>
  );
}
