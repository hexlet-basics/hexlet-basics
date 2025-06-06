import * as Routes from "@/routes.js";
import { Alert, Col, Container, Nav, Row, Tab } from "react-bootstrap";

import Chat from "@/components/Chat.tsx";
import MarkdownViewer from "@/components/MarkdownViewer.tsx";
import XssContent from "@/components/XssContent.tsx";
import { XBreadcrumb } from "@/components/breadcrumbs.tsx";
import LessonLayout from "@/pages/layouts/LessonLayout.tsx";
import type { BreadcrumbItem } from "@/types/index.ts";
import { Link, usePage } from "@inertiajs/react";
import i18next, { t } from "i18next";
import React, { useState } from "react";
import { useTranslation } from "react-i18next";
import App from "./components/App.tsx";
import type { LessonSharedProps } from "./types.ts";
import ContactMethodRequestingBlock from "@/pages/layouts/blocks/ContactMethodRequestingBlock.tsx";
import { useLessonStore } from "./store.tsx";

export default function Index() {
  const {
    courseCategory,
    landingPage,
    lessons,
    course,
    lessonMember,
    lesson,
    shouldAddContactMethod,
    canCreateAssistantMessage,
    previousMessages,
    // auth: { user },
  } = usePage<LessonSharedProps>().props;

  const { t: tCommon } = useTranslation("common");
  const { t: tViews } = useTranslation();

  const commonQuestions = t("languages.lessons.show.common_questions", {
    returnObjects: true,
  });

  const [focusesCount, setFocusCount] = useState(0);

  const handleSelect = (selectedKey: string | null) => {
    if (selectedKey == "assistant") {
      setFocusCount((state) => state + 1)
    }
  }

  const items: BreadcrumbItem[] = [
    {
      name: courseCategory?.name ?? '-',
      url: courseCategory ? Routes.language_category_url(courseCategory.slug!) : '#',
    },
    {
      name: landingPage.header!,
      url: Routes.language_url(landingPage.slug!),
    },
    {
      name: lesson.name!,
      url: Routes.language_lesson_url(landingPage.slug!, lesson.slug!),
    },
  ];

  const userCode = useLessonStore((state) => state.content);
  const output = useLessonStore((state) => state.output);

  return (
    <LessonLayout>
      <Container fluid className="overflow-hidden mb-1 x-h-md-100">
        <Row className="x-h-md-100">
          <Col className="x-h-md-100 col-12 col-md-6 col-lg-5 mb-3 mb-md-0 position-relative border-end">
            <Tab.Container id="left-tabs-example" defaultActiveKey="lesson">
              <div className="x-h-md-100 d-flex flex-column">
                <Nav variant="underline" onSelect={handleSelect} fill justify className="mb-3 small">
                  <Nav.Item>
                    <Nav.Link className="link-body-emphasis" eventKey="lesson">
                      {t("languages.lessons.show.lesson")}
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link className="link-body-emphasis" eventKey="assistant">
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

                <Tab.Content className="x-h-md-100 overflow-hidden">
                  <Tab.Pane
                    eventKey="lesson"
                    className="overflow-auto x-h-md-100"
                  >
                    <XBreadcrumb className="small" items={items} />

                    {/* {user.guest && ( */}
                    {/*   <Alert */}
                    {/*     variant="info" */}
                    {/*     className="border-0 small text-center" */}
                    {/*   > */}
                    {/*     <XssContent> */}
                    {/*       {t( */}
                    {/*         "languages.lessons.show.sign_up_for_tracking_progress_html", */}
                    {/*         { */}
                    {/*           link: Routes.new_user_path(), */}
                    {/*         }, */}
                    {/*       )} */}
                    {/*     </XssContent> */}
                    {/*   </Alert> */}
                    {/* )} */}

                    <div className="hexlet-basics-content">
                      <h1 className="h2">{`${landingPage.header}: ${lesson.name}`}</h1>

                      {shouldAddContactMethod && (
                        <div className="mt-3"><ContactMethodRequestingBlock /></div>
                      )}

                      <MarkdownViewer allowHtml>{lesson.theory || ""}</MarkdownViewer>
                      <h2 className="h3">
                        {t("languages.lessons.show.instructions")}
                      </h2>
                      <MarkdownViewer allowHtml>
                        {lesson.instructions || ""}
                      </MarkdownViewer>
                    </div>

                    {lesson.tips.length > 0 && (
                      <div>
                        <h2 className="h3">
                          {t("languages.lessons.show.tips")}
                        </h2>
                        <ul>
                          {lesson.tips.map((t) => (
                            <li key={t}>
                              <MarkdownViewer allowHtml>{t}</MarkdownViewer>
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

                    {course.hexlet_program_landing_page && (
                      <Alert variant="primary" className="my-4 small text-center shadow-sm border-0">
                        <a target="_blank" href={`${course.hexlet_program_landing_page}?utm_source=code-basics&utm_medium=referral`}>
                          {tViews('languages.lessons.show.profession_description')}
                        </a>
                      </Alert>
                    )}

                    <div className="my-4">
                      {commonQuestions.map((v) => (
                        <details
                          key={v.question}
                          className="mt-1 border rounded"
                        >
                          <summary className="p-2">{v.question}</summary>
                          <div className="px-2 pt-2">
                            <MarkdownViewer allowHtml>{v.answer}</MarkdownViewer>
                          </div>
                        </details>
                      ))}
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
                  <Tab.Pane
                    eventKey="assistant"
                    className="overflow-auto x-h-md-100"
                  >
                    {i18next.language === "ru" && (
                      <Alert className="small">
                        <XssContent>
                          {t("languages.lessons.show.if_stuck_html", {
                            url: tCommon("community_url"),
                          })}
                        </XssContent>
                      </Alert>
                    )}
                    <Chat
                      focusesCount={focusesCount}
                      previousMessages={previousMessages}
                      enabled={canCreateAssistantMessage}
                      userCode={userCode}
                      output={output}
                      course={course}
                      lesson={lesson}
                      lessonMember={lessonMember}
                    />
                  </Tab.Pane>
                  <Tab.Pane
                    eventKey="navigation"
                    className="overflow-auto x-h-md-100"
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
                              landingPage.language.slug!,
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

          <Col className="x-h-md-100 col-12 col-md-6 col-lg-7 mb-3 mb-md-0 position-relative col">
            <Tab.Container id="left-tabs-example" defaultActiveKey="editor">
              <div className="x-h-md-100 d-flex flex-column">
                <App />
              </div>
            </Tab.Container>
          </Col>
        </Row>
      </Container>
    </LessonLayout>
  );
}
