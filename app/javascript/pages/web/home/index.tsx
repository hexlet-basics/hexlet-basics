import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { TypeAnimation } from "react-type-animation";

import * as Routes from "@/routes.js";
import type {
  BlogPost,
  LanguageMember,
  Review,
  User,
} from "@/types/serializers";

import BlogPostBlock from "@/components/BlogPostBlock";
import CourseBlock from "@/components/CourseBlock";
import SignUpFormBlock from "@/components/SignUpFormBlock";
import XssContent from "@/components/XssContent";
import { getImageUrl } from "@/images";
import codeImagePathEn from "@/images/code-basics-coding-en.png";
import codeImagePathRu from "@/images/code-basics-coding-ru.png";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { LeadCrud, SharedProps } from "@/types";
import { Head, usePage } from "@inertiajs/react";
import { Accordion, Col, Container, Row } from "react-bootstrap";
import MarkdownViewer from "@/components/MarkdownViewer";
import type { Question, FAQPage, WithContext } from "schema-dts";
import LeadFormBlock from "@/components/LeadFormBlock";
import i18next from "i18next";

type Props = PropsWithChildren & {
  blogPosts: BlogPost[];
  courseMembersByCourseId: LanguageMember[];
  reviews: Review[];
  newUser: User;
  lead: LeadCrud;
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

// TODO: move to locales
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

export default function Index({
  blogPosts,
  newUser,
  lead,
  courseMembersByCourseId,
}: Props) {
  const { t } = useTranslation();
  // FIXME: add en version
  const { t: tFaq } = useTranslation("faq");

  const {
    locale,
    landingPagesForLists,
    auth: { user },
  } = usePage<SharedProps>().props;

  const faq = tFaq("main", { returnObjects: true });

  const entities: Question[] = Object.values(faq).map((item) => {
    return {
      "@type": "Question",
      "name": item.question,
      "acceptedAnswer": {
        "@type": "Answer",
        "text": item.answer,
      }
    }
  })
  // https://developers.google.com/search/docs/appearance/structured-data/faqpage
  const qaSchema: WithContext<FAQPage> = {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    "mainEntity": entities,
  };

  return (
    <ApplicationLayout>
      <Head>
        {Object.keys(faq).length > 0 && <script type="application/ld+json">{JSON.stringify(qaSchema)}</script>}
      </Head>
      <Container className="mb-lg-5 py-5">
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
              <div className="lead mb-3">
                <XssContent>
                  {t("home.hero.fastest_way_to_start_coding")}
                </XssContent>
              </div>
              <div className="d-grid gap-2 d-md-flex justify-content-md-start mb-4 mb-lg-3">
                <a
                  className="btn btn-primary btn-lg px-4 me-md-2 fw-bold"
                  href="#courses"
                >
                  {t("home.hero.try")}
                </a>
              </div>
            </div>
            <div className="col-lg-4 offset-lg-1 p-0 overflow-hidden shadow-lg text-center">
              <img
                alt="Code Basics learning preview"
                className="rounded-lg-3"
                height="356"
                src={codeImagePaths[locale]}
                width="720"
              />
            </div>
          </div>
        </div>
      </Container>

      <div className="mb-5">
        <Container className="py-4">
          <h2 className="mb-4">
            <a
              id="courses"
              className="text-decoration-none link-body-emphasis"
              href="#courses"
            >
              {t("home.languages.courses")}
            </a>
          </h2>

          <div className="row row-cols-2 row-cols-md-3 row-cols-lg-4">
            {landingPagesForLists.map((lp) => (
              <div className="col mb-3" key={lp.id}>
                <CourseBlock
                  landingPage={lp}
                  courseMember={courseMembersByCourseId[lp.id]}
                />
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
                <Accordion.Header as="h3">{value.question}</Accordion.Header>
                <Accordion.Body>
                  <MarkdownViewer>{value.answer}</MarkdownViewer>
                </Accordion.Body>
              </Accordion.Item>
            ))}
          </Accordion>
        </Container>
      </div>

      {user.guest && (
        <Container>
          <Row className="align-items-center py-5">
            <Col className="col-12 col-lg-7 fw-bold display-4 mb-5">
              {t("home.index.join")}
            </Col>
            <Col className="col-12 col-lg-5 mx-auto d-flex flex-column">
              <div className="bg-body-tertiary p-4 p-md-5 border rounded-3">
                <SignUpFormBlock user={newUser} />
              </div>
            </Col>
          </Row>
        </Container>
      )}

      {!user.guest && i18next.language === 'ru' && (
        <Container>
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
        </Container>
      )}

    </ApplicationLayout>
  );
}
