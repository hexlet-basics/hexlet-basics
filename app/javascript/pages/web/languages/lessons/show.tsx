import type { PropsWithChildren } from "react";
import { Alert, Col, Container, Nav, Row, Tab } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import XssContent from "@/components/XssContent";
import { XBreadcrumb } from "@/components/breadcrumbs";
import LessonLayout from "@/pages/layouts/LessonLayout";
import * as Routes from "@/routes.js";
import type { BreadcrumbItem } from "@/types/types";
import { Link, usePage } from "@inertiajs/react";
import React from "react";
import Markdown from "react-markdown";
import rehypeHighlight from "rehype-highlight";
import rehypeRaw from "rehype-raw";
import EditorBlock from "./show/index.tsx";
import type { Props } from "./show/types";

const rehypePlugins = [rehypeHighlight, rehypeRaw];

export default function Show() {
  const {
    courseCategory,
    course,
    lessons,
    lesson,
    prevLesson,
    nextLesson,
    auth: { user },
  } = usePage<Props>().props;
  const { t } = useTranslation();

  const commonQuestions = t("languages.lessons.show.common_questions", {
    returnObjects: true,
  });

  const items: BreadcrumbItem[] = [
    {
      name: courseCategory.name!,
      url: Routes.language_category_path(courseCategory.slug!),
    },
    {
      name: course.name!,
      url: Routes.language_path(course.slug!),
    },
    {
      name: lesson.name!,
      url: Routes.language_lesson_path(course.slug!, lesson.slug!),
    },
  ];

  return (
    <LessonLayout>
      <Container fluid className="overflow-hidden mb-1 h-100">
        <Row className="h-100">
          <Col className="h-100 col-12 col-md-6 col-lg-5 mb-3 mb-md-0 position-relative border-end">
            <Tab.Container id="left-tabs-example" defaultActiveKey="lesson">
              <div className="h-100 d-flex flex-column">
                <Nav variant="underline" fill justify className="mb-3 small">
                  <Nav.Item>
                    <Nav.Link className="link-body-emphasis" eventKey="lesson">
                      {t("languages.lessons.show.lesson")}
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link className="link-body-emphasis" eventKey="discuss">
                      {t("languages.lessons.show.discuss")}
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link
                      className="link-body-emphasis"
                      eventKey="navigation"
                    >
                      {t("languages.lessons.show.navigation")}
                    </Nav.Link>
                  </Nav.Item>
                </Nav>

                <Tab.Content className="h-100 overflow-hidden">
                  <Tab.Pane eventKey="lesson" className="overflow-auto h-100">
                    <XBreadcrumb className="small" items={items} />

                    {user.guest && (
                      <Alert variant="info" className="border-0">
                        <div
                          // biome-ignore lint/security/noDangerouslySetInnerHtml: <explanation>
                          dangerouslySetInnerHTML={{
                            __html: t(
                              "languages.lessons.show.sign_up_for_tracking_progress_html",
                              {
                                name: course.name,
                                link: Routes.new_user_path(),
                              },
                            ),
                          }}
                        />
                      </Alert>
                    )}

                    <h1 className="h2">{`${course.name}: ${lesson.name}`}</h1>
                    <Markdown
                      className="hexlet-basics-content"
                      rehypePlugins={rehypePlugins}
                    >
                      {lesson.theory}
                    </Markdown>
                    <h2 className="h3">
                      {t("languages.lessons.show.instructions")}
                    </h2>
                    <Markdown rehypePlugins={rehypePlugins}>
                      {lesson.instructions}
                    </Markdown>

                    {lesson.tips.length > 0 && (
                      <div>
                        <h2 className="h3">
                          {t("languages.lessons.show.tips")}
                        </h2>
                        <ul>
                          {lesson.tips.map((t) => (
                            <li key={t}>
                              <Markdown>{t}</Markdown>
                            </li>
                          ))}
                        </ul>
                      </div>
                    )}

                    {lesson.definitions.length > 0 && (
                      <div>
                        <h2 className="h3">
                          {t("languages.lessons.show.definitions")}
                        </h2>
                        <dl>
                          {lesson.definitions.map((d) => (
                            <React.Fragment key={d.name}>
                              <dt>{d.name}</dt>
                              <dd>{d.description}</dd>
                            </React.Fragment>
                          ))}
                        </dl>
                      </div>
                    )}

                    <div className="my-4">
                      {commonQuestions.map((v) => (
                        <details
                          key={v.question}
                          className="mt-1 border rounded"
                        >
                          <summary className="p-2">{v.question}</summary>
                          <Markdown className="px-2 pt-2">{v.answer}</Markdown>
                        </details>
                      ))}
                    </div>

                    <div className="mb-3 d-flex justify-content-center">
                      {prevLesson && (
                        <Link
                          href={Routes.language_lesson_path(
                            course.slug!,
                            prevLesson.slug!,
                          )}
                          className="text-decoration-none link-body-emphasis me-2"
                        >
                          {t("languages.lessons.show.prev")}
                        </Link>
                      )}
                      {nextLesson && (
                        <Link
                          className="text-decoration-none link-body-emphasis"
                          href={Routes.language_lesson_path(
                            course.slug!,
                            nextLesson.slug!,
                          )}
                        >
                          {t("languages.lessons.show.next")}
                        </Link>
                      )}
                    </div>

                    <div className="small text-muted py-2">
                      <span className="me-2">
                        {t("languages.lessons.show.issues")}
                      </span>
                      <a
                        href={lesson.source_code_url!}
                        target="_blank"
                        rel="noreferrer"
                        className="link-body-emphasis"
                      >
                        <i className="bi bi-github" />
                      </a>
                    </div>
                  </Tab.Pane>
                  <Tab.Pane eventKey="discuss">
                    <XssContent>
                      {t("languages.lessons.show.if_stuck_html")}
                    </XssContent>
                  </Tab.Pane>
                  <Tab.Pane
                    eventKey="navigation"
                    className="overflow-auto h-100"
                  >
                    <ul className="list-unstyled">
                      {lessons.map((l) => (
                        <li key={l.id} className="mb-1">
                          {l.natural_order}
                          {". "}
                          <Link
                            as={l.slug === lesson.slug ? "b" : "a"}
                            className="link-body-emphasis"
                            href={Routes.language_lesson_path(
                              course.slug!,
                              l.slug!,
                            )}
                          >
                            {l.name}
                          </Link>
                        </li>
                      ))}
                    </ul>
                  </Tab.Pane>
                </Tab.Content>
              </div>
            </Tab.Container>
          </Col>

          <Col className="h-100 col-12 col-md-6 col-lg-7 mb-3 mb-md-0 position-relative col">
            <Tab.Container id="left-tabs-example" defaultActiveKey="editor">
              <div className="h-100 d-flex flex-column">
                <EditorBlock />
                {/* <Nav variant="underline" fill justify className="mb-3 small"> */}
                {/*   <Nav.Item> */}
                {/*     <Nav.Link className="link-body-emphasis" eventKey="editor"> */}
                {/*       {t("languages.lessons.show.editor")} */}
                {/*     </Nav.Link> */}
                {/*   </Nav.Item> */}
                {/*   <Nav.Item> */}
                {/*     <Nav.Link className="link-body-emphasis" eventKey="output"> */}
                {/*       {t("languages.lessons.show.output")} */}
                {/*     </Nav.Link> */}
                {/*   </Nav.Item> */}
                {/*   <Nav.Item> */}
                {/*     <Nav.Link className="link-body-emphasis" eventKey="Tests"> */}
                {/*       {t("languages.lessons.show.tests")} */}
                {/*     </Nav.Link> */}
                {/*   </Nav.Item> */}
                {/*   <Nav.Item> */}
                {/*     <Nav.Link */}
                {/*       className="link-body-emphasis" */}
                {/*       eventKey="solution" */}
                {/*     > */}
                {/*       {t("languages.lessons.show.solution")} */}
                {/*     </Nav.Link> */}
                {/*   </Nav.Item> */}
                {/* </Nav> */}

                {/* <Tab.Content className="h-100 overflow-hidden p-2"> */}
                {/*   <Tab.Pane eventKey="lesson" className="overflow-auto h-100"> */}
                {/*     <h1 className="h2"> */}
                {/*       {course.name} */}
                {/*       {lesson.name} */}
                {/*     </h1> */}
                {/*     <Markdown rehypePlugins={rehypePlugins}> */}
                {/*       {lesson.theory} */}
                {/*     </Markdown> */}
                {/*     <h2 className="h3"> */}
                {/*       {t("languages.lessons.show.instructions")} */}
                {/*     </h2> */}
                {/*     <Markdown rehypePlugins={rehypePlugins}> */}
                {/*       {lesson.instructions} */}
                {/*     </Markdown> */}
                {/*   </Tab.Pane> */}
                {/*   <Tab.Pane eventKey="discuss">Second tab content</Tab.Pane> */}
                {/*   <Tab.Pane eventKey="navigation">Second tab content</Tab.Pane> */}
                {/* </Tab.Content> */}
              </div>
            </Tab.Container>
          </Col>
        </Row>
      </Container>
    </LessonLayout>
  );
}
