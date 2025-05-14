import { getImageUrl } from "@/images";
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
  const { t: tLinks } = useTranslation("links");

  const programmingBasicList = t(
    "for_school_teachers.programming_basic_list",
    { returnObjects: true },
  );

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
              href={user.guest ? Routes.new_user_path() : tLinks("hexlet_courses")}
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
            <p className="fs-5 text-body-secondary lh-sm pe-lg-4">
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
                {programmingBasicList.map((item, index) => (
                  <li key={index.toString()} className="d-flex gap-3 mb-3">
                    <span>–</span>
                    <span>{item}</span>
                  </li>
                  ))}
              </ul>
            </div>
            <div className="pe-lg-5">
              <div className="fs-5 fw-medium lh-sm">
                {t("for_school_teachers.interactive_approach")}
              </div>
              <ul className="text-body-secondary lh-1 list-unstyled my-4 pb-2 py-lg-2">
                {interactiveApproachList.map((item, index) => (
                  <li key={index.toString()} className="d-flex gap-3 mb-3">
                    <span>–</span>
                    <span>{item}</span>
                  </li>
                ))}
              </ul>
            </div>
            <div className="pe-lg-5">
              <div className="fs-5 fw-medium lh-sm">
                {t("for_school_teachers.early_career_guidance")}
              </div>
              <ul className="text-body-secondary lh-1 list-unstyled my-4 py-lg-2">
                {earlyCareerGuidanceList.map((item, index) => (
                  <li key={index.toString()} className="d-flex gap-3 mb-3">
                    <span>–</span>
                    <span>{item}</span>
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
                  href={tLinks('hexlet_courses')}
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
            {howToLearnCards.map((item, index) => (
              <Col key={index.toString()}>
                <div className="border rounded-4 p-3 p-lg-4 h-100 d-flex flex-column gap-4">
                  <img
                    src={getImageUrl(`for-school-teachers-page/${item.img}.svg`)}
                    width="90"
                    height="65"
                    alt={item.img}
                  />
                  <div className="lh-sm">
                    {item.title}
                  </div>
                  <div className="lh-sm">
                    {item.subtitle}
                  </div>
                </div>
              </Col>
              ))}
          </Row>
        </div>
        <div className="pt-5 pb-3 pb-lg-5">
          <Row className="bg-dark text-light rounded-4 d-flex align-items-end justify-content-between p-2 p-md-3 p-lg-4 m-0 gy-2">
            <Col className="col-md-10 col-lg-8 col-xl-7">
              <div className="h1 lh-1 fw-normal pe-lg-4 mb-4">
                {t("for_school_teachers.sign_up_and_start_learning")}
              </div>
            </Col>
            <Col className="col-12 col-lg-auto">
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
                  href={tLinks("hexlet_courses")}
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
