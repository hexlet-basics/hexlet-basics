import { type PropsWithChildren, useEffect } from "react";
import { Card, Col, Container, Nav, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import type { SharedProps } from "@/types/index.js";
import { Head, Link, usePage } from "@inertiajs/react";
import i18next from "i18next";
import NavbarBlock from "./NavbarBlock.tsx";
import XFlash from "@/components/XFlash.tsx";

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
    <>
      {/* <Head> */}
      {/*   <title>{header}</title> */}
      {/* </Head> */}
      <Container fluid className="mb-2">
        <NavbarBlock />
        <XFlash />
      </Container>
      <Container fluid className="h-100">
        <Row className="h-100">
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
                  href={Routes.admin_languages_path()}
                >
                  <i className="bi bi-terminal me-2" />
                  {tLayouts("web.admin.application.languages")}
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
            </Nav>
          </Col>
          <Col className="col-10">
            <h1 className="mb-4">{header}</h1>
            {children}
          </Col>
        </Row>
      </Container>
    </>
  );
}
