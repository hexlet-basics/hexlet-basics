import type { PropsWithChildren } from "react";
import { Col, Container, Nav, Row, Tab } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import { XBreadcrumb } from "@/components/breadcrumbs";
import LessonLayout from "@/pages/layouts/LessonLayout";
import * as Routes from "@/routes.js";
import type {
  Language,
  LanguageCategory,
  LanguageLesson,
} from "@/types/serializers";
import type { BreadcrumbItem, SharedProps } from "@/types/types";
import { usePage } from "@inertiajs/react";
import Markdown from "react-markdown";
import rehypeHighlight from "rehype-highlight";
import rehypeRaw from "rehype-raw";

type Props = PropsWithChildren & {
  courseCategory: LanguageCategory;
  course: Language;
  lesson: LanguageLesson;
};

const rehypePlugins = [rehypeHighlight, rehypeRaw];

export default function Show({ courseCategory, course, lesson }: Props) {
  const { suffix } = usePage<SharedProps>().props;
  const { t } = useTranslation();

  const items: BreadcrumbItem[] = [
    {
      name: courseCategory.name!,
      url: Routes.language_category_path(courseCategory.slug!, { suffix }),
    },
    {
      name: course.name!,
      url: Routes.language_path(course.slug!, { suffix }),
    },
    {
      name: lesson.name!,
      url: Routes.language_lesson_path(course.slug!, lesson.slug!, { suffix }),
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

                    {/* - unless signed_in? */}
                    {/*   .alert.alert-info.border-0.rounded-0.small */}
                    {/*     = t('.sign_up_for_tracking_progress_html', name: resource_language, link: new_user_path) */}

                    <h1 className="h2">
                      {`${course.name}: ${lesson.name}`}
                    </h1>
                    <Markdown rehypePlugins={rehypePlugins}>
                      {lesson.theory}
                    </Markdown>
                    <h2 className="h3">
                      {t("languages.lessons.show.instructions")}
                    </h2>
                    <Markdown rehypePlugins={rehypePlugins}>
                      {lesson.instructions}
                    </Markdown>

                    {/* .my-4 */}
                    {/*   - t('.common_questions').each do |data| */}
                    {/*     details.mt-1.border.rounded */}
                    {/*       summary.p-2 = data[:question] */}
                    {/*       .px-2.pt-2 == markdown2html data[:answer] */}

                    {/* - if @info.tips.any? */}
                    {/*   h2.h5.font-weight-bold.mt-4 = t('.tips') */}
                    {/*   ul.ps-4 */}
                    {/*     - @info.tips.each do |tip| */}
                    {/*       li == markdown2html(tip, options) */}
                    {/* - if @info.definitions.any? */}
                    {/*   h2.h5.font-weight-bold = t('.definitions') */}
                    {/*   ul.ps-4 */}
                    {/*     - @info.definitions.each do |definition| */}
                    {/*       li */}
                    {/*         / NOTE: add different separators for different locales */}
                    {/*         p == markdown2html("#{definition['name']}#{t('.separator')}#{definition['description']}") */}

                    {/*   .mb-3.d-flex.justify-content-center */}
                    {/*     .me-4 */}
                    {/*       - if @lesson_version.prev_lesson */}
                    {/*         = link_to language_lesson_path(resource_language.slug, @lesson_version.prev_lesson.slug), class: 'd-inline-block text-black text-decoration-none' do */}
                    {/*           span.me-2 ← */}
                    {/*           span = t('.prev') */}
                    {/*     - if @lesson_version.next_lesson */}
                    {/*       = link_to language_lesson_path(resource_language.slug, @lesson_version.next_lesson.slug), class: 'd-inline-block text-black text-decoration-none' do */}
                    {/*         span.me-2 = t('.next') */}
                    {/*         span → */}
                    {/**/}
                    {/* hr.my-4 */}
                    {/**/}
                    {/* .small.text-muted */}
                    {/*   = t('.issues') */}
                    {/*   ' */}
                    {/*   = link_to ExternalLinks.source_code_curl, get_lesson_source_code(@lesson_version, @info), target: '_blank', rel: 'nofollow noopener' */}
                  </Tab.Pane>
                  <Tab.Pane eventKey="discuss">Second tab content</Tab.Pane>
                  <Tab.Pane eventKey="navigation">Second tab content</Tab.Pane>
                </Tab.Content>
              </div>
            </Tab.Container>
          </Col>
          <Col className="h-100 col-12 col-md-6 col-lg-7 pl-md-0">
            <Tab.Container id="left-tabs-example" defaultActiveKey="editor">
              <div className="h-100">
                <Nav variant="underline" fill justify className="mb-3 small">
                  <Nav.Item>
                    <Nav.Link className="link-body-emphasis" eventKey="editor">
                      {t("languages.lessons.show.editor")}
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link className="link-body-emphasis" eventKey="output">
                      {t("languages.lessons.show.output")}
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link className="link-body-emphasis" eventKey="Tests">
                      {t("languages.lessons.show.tests")}
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link
                      className="link-body-emphasis"
                      eventKey="solution"
                    >
                      {t("languages.lessons.show.solution")}
                    </Nav.Link>
                  </Nav.Item>
                </Nav>

                <Tab.Content className="h-100 overflow-hidden p-2">
                  <Tab.Pane eventKey="lesson" className="overflow-auto h-100">
                    <XBreadcrumb items={items} />
                    <h1 className="h2">
                      {course.name}
                      {lesson.name}
                    </h1>
                    <Markdown rehypePlugins={rehypePlugins}>
                      {lesson.theory}
                    </Markdown>
                    <h2 className="h3">
                      {t("languages.lessons.instructions")}
                    </h2>
                    <Markdown rehypePlugins={rehypePlugins}>
                      {lesson.instructions}
                    </Markdown>
                  </Tab.Pane>
                  <Tab.Pane eventKey="discuss">Second tab content</Tab.Pane>
                  <Tab.Pane eventKey="navigation">Second tab content</Tab.Pane>
                </Tab.Content>
              </div>
            </Tab.Container>
          </Col>
        </Row>
      </Container>
    </LessonLayout>
  );
}
