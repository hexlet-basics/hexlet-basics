import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import Markdown from "react-markdown";
import { TypeAnimation } from "react-type-animation";

import * as Routes from "@/routes.js";
import type {
  BlogPost,
  LanguageCategory,
  Review,
  User,
} from "@/types/serializers";

import BlogPostBlock from "@/components/BlogPostBlock";
import CourseBlock from "@/components/CourseBlock";
import { XForm, XInput } from "@/components/forms";
import { getImageUrl } from "@/images";
import codeImagePathEn from "@/images/code-basics-coding-en.png";
import codeImagePathRu from "@/images/code-basics-coding-ru.png";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { SharedProps } from "@/types";
import { Link, usePage } from "@inertiajs/react";
import { Accordion, Button, Col, Container, Form, Row } from "react-bootstrap";
import { Submit } from "use-inertia-form";
import XssContent from "@/components/XssContent";

type Props = PropsWithChildren & {
  blogPosts: BlogPost[];
  courseCategories: LanguageCategory[];
  reviews: Review[];
  newUser: User;
};

const sequence = [
  "TypeScript",
  1000,
  "Java",
  1000,
  "Python",
  1000,
  "PHP",
  1000,
  "Ruby",
  1000,
  "HTML",
  1000,
  "CSS",
  1000,
  "Go",
  1000,
];

const codeImagePaths = {
  ru: codeImagePathRu,
  en: codeImagePathEn,
};

const reviews = {
  ru: [
    {
      name: "Александр Авдошкин",
      avatar: getImageUrl("avdoshkin.jpg"),
      body: `Если бы не коронавирус, выполнил бы всё в заход (в смысле каждый день по несколько пунктов в теме).
              Изучаю с нуля, ваш портал очень ориентирован на новичков. Спасибо вам большое!`,
    },
    {
      name: "Сергей Тюрин",
      avatar: getImageUrl("tyrin.jpg"),
      body: `Очень всё доступно даже для полного профана вроде меня. Эта вводная по JS вошла в мой туговатый ум,
              складно как недостающий пазл. Всем кидаю линк на эту страничку.`,
    },
    {
      name: "Элиях Клейман",
      avatar: getImageUrl("user-avatar.png"),
      body: `Для меня это первый курс для новичка. Понравилось тем, что вся информация структурирована
              и дана по мере изучения материала в иерархичном порядке, что значительно повышает и желание к обучению`,
    },
  ],
  en: [
    {
      name: "Aleksandr Avdoshkin",
      avatar: getImageUrl("avdoshkin.jpg"),
      body: "As someone with zero coding skills, I'd say that CodeBasics is focused on newcomers. Thank you very much!",
    },
    {
      name: "Sergei Tyurin",
      avatar: getImageUrl("tyrin.jpg"),
      body: "This is all very approachable even for a dummy like me. Now I show people this platform when I get the chance.",
    },
    {
      name: "Eliyah Kleyman",
      avatar: getImageUrl("user-avatar.png"),
      body: `For me, it was my very first programming course. I liked it because all the information is very well 
        structured and given in a clear hierarchical order. It motivated me a lot to move forward in my studies.`,
    },
  ],
};

export default function Index({ blogPosts, newUser, courseCategories }: Props) {
  const { t } = useTranslation();
  // FIXME: en version
  const { t: tFaq } = useTranslation("faq");
  const { t: tHelpers } = useTranslation("helpers");

  const {
    locale,
    courses,
    auth: { user },
  } = usePage<SharedProps>().props;

  const faq = tFaq("main", { returnObjects: true });

  return (
    <ApplicationLayout>
      <Container className="mb-5 py-5">
        <div className="bg-body-tertiary p-4 pb-0 pt-lg-5 align-items-center border shadow-sm rounded-3">
          <div className="row">
            <div className="col-lg-7 p-3 p-lg-5 pt-lg-3">
              <h1 className="h6 text-muted">
                {t("home.hero.free_programming_courses")}
              </h1>
              <div className="display-4 fw-bold lh-1 mb-3">
                {t("home.hero.learn_html")}{" "}
                <TypeAnimation
                  className="text-primary"
                  sequence={sequence}
                  wrapper="span"
                  speed={5}
                  repeat={Number.POSITIVE_INFINITY}
                />
              </div>
              <p className="lead">
                {t("home.hero.fastest_way_to_start_coding")}
              </p>
              <div className="d-grid gap-2 d-md-flex justify-content-md-start mb-4 mb-lg-3">
                <a
                  className="btn btn-primary btn-lg px-4 me-md-2 fw-bold"
                  href="#courses"
                >
                  Попробовать
                </a>
              </div>
            </div>
            <div className="col-lg-4 offset-lg-1 p-0 overflow-hidden shadow-lg text-center">
              <img
                alt="Как работает обучение на code-basics"
                className="rounded-lg-3"
                height="356"
                src={codeImagePaths[locale]}
                width="720"
              />
            </div>
          </div>
        </div>
      </Container>

      <div className="py-5 mb-5">
        <Container className="py-4">
          <div className="d-flex flex-column flex-sm-row align-items-sm-center mb-4">
            <h2 className="me-auto mb-3 mb-sm-0">
              <a
                id="courses"
                className="text-decoration-none link-body-emphasis"
                href="#courses"
              >
                {t("home.languages.courses")}
              </a>
            </h2>

            <div className="d-flex flex-wrap">
              {courseCategories.map((c) => (
                <a
                  key={c.id}
                  href={Routes.language_category_path(c.slug!)}
                  className="mb-2 mb-sm-0 me-2 me-sm-0 ms-sm-2 fw-light text-decoration-none fs-5 badge text-bg-light p-2 p-sm-3 rounded-pill border"
                >
                  {c.name}
                </a>
              ))}
            </div>
          </div>

          <div className="row row-cols-2 row-cols-md-3 row-cols-lg-4">
            {courses.map((course) => (
              <div className="col mb-3" key={course.id}>
                <CourseBlock course={course} />
              </div>
            ))}
          </div>
        </Container>
      </div>

      <div className="bg-body-tertiary py-5 mb-5">
        <Container className="mb-5">
          <div className="d-flex">
            <h2 className="me-auto mt-auto">{t("home.index.reviews")}</h2>
            <div className="mt-auto">
              <a
                href={Routes.reviews_path()}
                className="text-decoration-none text-muted small"
              >
                {t("home.index.all_reviews")}
              </a>
            </div>
          </div>
          <hr className="mb-5 mt-1" />

          <div className="row g-4 row-cols-1 row-cols-sm-2 row-cols-lg-3">
            {reviews[locale].map((review) => (
              <div key={review.name} className="col">
                <div className="d-flex mb-3">
                  <img
                    src={review.avatar}
                    className="rounded-circle flex-shrink-0"
                    width={50}
                    height={50}
                    alt={`Аватар пользователя ${review.name}`}
                  />
                  <div className="ms-3">
                    <div className="fw-bold">{review.name}</div>
                    <div className="small">
                      {/* {t(".course_html", { */}
                      {/*   link: ( */}
                      {/*     <a */}
                      {/*       href={language_path(jsCourse.slug)} */}
                      {/*       className="text-dark" */}
                      {/*     > */}
                      {/*       {jsCourse} */}
                      {/*     </a> */}
                      {/*   ), */}
                      {/* })} */}
                    </div>
                  </div>
                </div>
                <div className="fst-italic">{review.body}</div>
              </div>
            ))}
          </div>
        </Container>
      </div>

      {blogPosts.length > 0 && (
        <Container className="mb-5">
          <>
            <div className="d-flex">
              <h2 className="me-auto mt-auto">{t("home.index.blog_posts")}</h2>
              <div className="mt-auto">
                <a
                  href={Routes.blog_posts_path()}
                  className="text-decoration-none text-muted small"
                >
                  {t("home.index.all_blog_posts")}
                </a>
              </div>
            </div>
            <hr className="mb-5 mt-1" />
            <Row className="row-cols-sm-2 row-cols-md-3 row-cols-1">
              {blogPosts.map((post) => (
                <Col key={post.id}>
                  <BlogPostBlock post={post} />
                </Col>
              ))}
            </Row>
          </>
        </Container>
      )}

      {locale === "ru" && (
        <div className="bg-body-tertiary mb-5 py-5">
          <Container className="mb-5">
            <h2>
              <a
                id="faq"
                className="text-decoration-none link-body-emphasis"
                href="#faq"
              >
                {tFaq("header")}
              </a>
            </h2>
            <hr className="mb-5" />
            <Accordion defaultActiveKey="0">
              {Object.entries(faq).map(([key, value], index) => (
                <Accordion.Item eventKey={String(index)} key={key}>
                  <Accordion.Header>{value.question}</Accordion.Header>
                  <Accordion.Body>
                    <Markdown>{value.answer}</Markdown>
                  </Accordion.Body>
                </Accordion.Item>
              ))}
            </Accordion>
          </Container>
        </div>
      )}

      {user.guest && (
        <div className="">
          <div className="container">
            <div className="row align-items-center g-lg-5 py-5">
              <div className="col-lg-7 fw-bold display-4">
                {t("home.index.join")}
                {/* {cache([I18n.locale], { expiresIn: "1d" }, () => ( */}
                {/*   <> */}
                {/*     <h1 className="display-4 fw-bold lh-1 mb-3"> */}
                {/*       {t(".join", { count: User.count })} */}
                {/*     </h1> */}
                {/*     Uncomment this if the paragraph is needed */}
                {/*     <p className="col-lg-10 fs-4">Forever free.</p> */}
                {/*   </> */}
                {/* ))} */}
              </div>
              <div className="col-md-10 mx-auto col-lg-5">
                <XForm
                  model="user_sign_up_form"
                  data={newUser}
                  to={Routes.users_path()}
                  className="d-flex flex-column bg-body-tertiary p-4 p-md-5 border rounded-3"
                >
                  <XInput name="first_name" autoComplete="name" />
                  <XInput name="email" autoComplete="email" />
                  <XInput name="password" type="password" autoComplete="current-password" />
                  <div className="text-end text-muted small mb-4">
                    {t("users.new.have_account")}{" "}
                    <Link
                      href={Routes.new_session_path()}
                      className="text-decoration-none"
                    >
                      {t("users.new.sign_in")}
                    </Link>
                  </div>
                  <Submit
                    // size="lg"
                    className="btn w-100 btn-lg btn-primary mb-3"
                    // type="submit"
                    // disabled={form.processing}
                  >
                    {tHelpers("submit.user_sign_up_form.create")}
                  </Submit>
                  <XssContent className="small text-muted">
                    {t("users.new.confirmation_html", {
                      url: Routes.page_path("tos"),
                    })}
                  </XssContent>
                </XForm>
              </div>
            </div>
          </div>
        </div>
      )}
    </ApplicationLayout>
  );
}
