import type { PropsWithChildren } from "react";
import { Col, Container, Nav, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import XFlash from "@/components/XFlash.tsx";
import { Link } from "@inertiajs/react";
import RootLayout from "./RootLayout.tsx";
import FooterBlock from "./blocks/FooterBlock.tsx";
import NavbarBlock from "./blocks/NavbarBlock.tsx";

type Props = PropsWithChildren & {
  header: string;
};

export default function AdminLayout({ children, header }: Props) {
  const { t: tLayouts } = useTranslation("layouts");
  const { t: tCommon } = useTranslation("common");

  // const { locale, suffix } = usePage<SharedProps>().props;
  //
  // useEffect(() => {
  //   i18next.changeLanguage(locale);
  //   Routes.configure({ default_url_options: { suffix } });
  // }, [suffix, locale]);

  return (
    <RootLayout>
      {/* <Head> */}
      {/*   <title>{header}</title> */}
      {/* </Head> */}
      <Container fluid className="mb-2">
        <NavbarBlock />
        <XFlash />
      </Container>
      <Container fluid className="h-100">
        <Row className="mb-5">
          <Col className="col-2 border-end h-100">
            <Nav className="flex-column">
              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_root_path()}
                >
                  <i className="bi bi-house me-2" />
                  {tLayouts("web.admin.application.dashboard")}
                </Link>
              </Nav.Item>
              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_language_categories_path()}
                >
                  <i className="bi bi-file-text me-2" />
                  {tLayouts("web.admin.application.language_categories")}
                </Link>
              </Nav.Item>
              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_languages_path()}
                >
                  <i className="bi bi-terminal me-2" />
                  {tLayouts("web.admin.application.languages")}
                </Link>
              </Nav.Item>
              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_language_landing_pages_path()}
                >
                  <i className="bi bi-terminal me-2" />
                  {tLayouts("web.admin.application.language_landing_pages")}
                </Link>
              </Nav.Item>
              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_reviews_path()}
                >
                  <i className="bi bi-chat-left-quote me-2" />
                  {tLayouts("web.admin.application.reviews")}
                </Link>
              </Nav.Item>
              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_blog_posts_path()}
                >
                  <i className="bi bi-file-text me-2" />
                  {tLayouts("web.admin.application.blog_posts")}
                </Link>
              </Nav.Item>
              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_messages_path()}
                >
                  <i className="bi bi-file-text me-2" />
                  {tLayouts("web.admin.application.messages")}
                </Link>
              </Nav.Item>
              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_language_lesson_members_path()}
                >
                  <i className="bi bi-file-text me-2" />
                  {tLayouts("web.admin.application.language_lesson_members")}
                </Link>
              </Nav.Item>
              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_surveys_path()}
                >
                  <i className="bi bi-patch-question me-2" />
                  {tLayouts("web.admin.application.surveys")}
                </Link>
              </Nav.Item>
              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_survey_answers_path()}
                >
                  <i className="bi bi-signpost-2 me-2" />
                  {tLayouts("web.admin.application.survey_answers")}
                </Link>
              </Nav.Item>

              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_leads_url()}
                >
                  <i className="bi bi-person-check me-2" />
                  {tLayouts("web.admin.application.leads")}
                </Link>
              </Nav.Item>

              <hr />

              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.admin_management_users_path()}
                >
                  <i className="bi bi-people me-2" />
                  {tLayouts("web.admin.application.users")}
                </Link>
              </Nav.Item>

              <hr />

              <Nav.Item>
                <Link
                  className="nav-link fw-medium link-body-emphasis"
                  href={Routes.surveys_admin_analytics_path()}
                >
                  <i className="bi bi-people me-2" />
                  {tLayouts("web.admin.application.surveys_analytics")}
                </Link>
              </Nav.Item>
            </Nav>
          </Col>
          <Col className="col-10">
            <h1 className="mb-4">{header}</h1>
            {children}
          </Col>
        </Row>
        <FooterBlock />
      </Container>
    </RootLayout>
  );
}
