import { type PropsWithChildren, useEffect } from "react";
import { Col, Container, Nav, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import XFlash from "@/components/XFlash.tsx";
import type { SharedProps } from "@/types/index.js";
import { Link, usePage } from "@inertiajs/react";
import i18next from "i18next";
import NavbarBlock from "./NavbarBlock.tsx";

type Props = PropsWithChildren & {
  header?: string;
};

export default function ApplicationLayout({ children, header }: Props) {
  const { t: tLayouts } = useTranslation("layouts");
  const { t: tCommon } = useTranslation("common");
  const { locale } = usePage<SharedProps>().props;

  return (
    <>
      <Container className="py-2">
        <NavbarBlock className="pb-3 border-bottom" />
        <XFlash />
      </Container>
      <Container className="mb-4">{header && <h1 className="mt-5">{header}</h1>}</Container>
      {children}
      <footer className="mt-5 bg-body-tertiary py-5">
        <Container>
          <Row>
            <Col>
              <Nav className="flex-column">
                <Nav.Item className="h5">+7 (495) 085 21 62</Nav.Item>
                <Nav.Item className="h5">8 800 100 22 47</Nav.Item>
                <Nav.Item className="mb-3">
                  <a href="mailto:support@hexlet.io">support@hexlet.io</a>
                </Nav.Item>
                <Nav.Item>{tCommon("address")}</Nav.Item>
              </Nav>
            </Col>
            <Col>
              <Nav className="flex-column mb-3">
                <Nav.Item>
                  <Link
                    className="link-body-emphasis text-decoration-none"
                    href={Routes.page_path("about")}
                  >
                    {tLayouts("shared.footer.about")}
                  </Link>
                </Nav.Item>
                <Nav.Item>
                  <Link
                    className="link-body-emphasis text-decoration-none"
                    href={Routes.blog_posts_path()}
                  >
                    {tLayouts("shared.footer.blog")}
                  </Link>
                </Nav.Item>
                <Nav.Item>
                  <Link
                    className="link-body-emphasis text-decoration-none"
                    href={Routes.reviews_path()}
                  >
                    {tLayouts("shared.footer.reviews")}
                  </Link>
                </Nav.Item>
                <Nav.Item>
                  <Link
                    className="link-body-emphasis text-decoration-none"
                    href={Routes.page_path("authors")}
                  >
                    {tLayouts("shared.footer.authors")}
                  </Link>
                </Nav.Item>
                <Nav.Item>
                  <a
                    className="link-body-emphasis text-decoration-none"
                    href="https://ttttt.me/hexlet_ru"
                  >
                    {tLayouts("shared.footer.community")}
                  </a>
                </Nav.Item>
              </Nav>
            </Col>
            <Col>
              <Nav className="flex-column">
                <Nav.Item>
                  <a
                    target="_blank"
                    rel="noreferrer"
                    className="link-body-emphasis text-decoration-none"
                    href="https://hexlet.io"
                  >
                    {tLayouts("shared.footer.hexlet")}
                  </a>
                </Nav.Item>
                <Nav.Item>
                  <a
                    target="_blank"
                    rel="noreferrer"
                    className="link-body-emphasis text-decoration-none"
                    href="https://cv.hexlet.io"
                  >
                    {tLayouts("shared.footer.hexlet-cv")}
                  </a>
                </Nav.Item>
                <Nav.Item>
                  <a
                    target="_blank"
                    rel="noreferrer"
                    className="link-body-emphasis text-decoration-none"
                    href="https://sicp.hexlet.io"
                  >
                    {tLayouts("shared.footer.hexlet-sicp")}
                  </a>
                </Nav.Item>
                <Nav.Item>
                  <a
                    target="_blank"
                    rel="noreferrer"
                    className="link-body-emphasis text-decoration-none"
                    href="https://friends.hexlet.io"
                  >
                    {tLayouts("shared.footer.hexlet-friends")}
                  </a>
                </Nav.Item>
              </Nav>
            </Col>
            <Col>
              <Nav className="flex-column">
                <Nav.Item>
                  <Link
                    className="link-body-emphasis text-decoration-none"
                    href={Routes.page_path("tos")}
                  >
                    {tLayouts("shared.footer.tos")}
                  </Link>
                </Nav.Item>
                <Nav.Item>
                  <Link
                    className="link-body-emphasis text-decoration-none"
                    href={Routes.page_path("privacy")}
                  >
                    {tLayouts("shared.footer.privacy")}
                  </Link>
                </Nav.Item>
              </Nav>
            </Col>
          </Row>
          <div className="d-flex flex-column flex-sm-row justify-content-between pt-4 my-4 border-top">
            <div>{`© ${new Date().getFullYear()} ${tCommon("legal_name")}`}</div>
            <ul className="fs-3 d-flex list-unstyled">
              <li className="me-3">
                <a
                  target="_blank"
                  rel="noreferrer"
                  className="link-body-emphasis"
                  href="https://github.com/hexlet-basics"
                >
                  <i className="bi bi-github" />
                </a>
              </li>
              {locale === "ru" && (
                <li className="me-3">
                  <a
                    target="_blank"
                    rel="noreferrer"
                    className="link-body-emphasis"
                    href="https://ttttt.me/hexlet_ru"
                  >
                    <i className="bi bi-telegram" />
                  </a>
                </li>
              )}
              <li>
                <a
                  target="_blank"
                  rel="noreferrer"
                  className="link-body-emphasis"
                  href="https://www.youtube.com/@HexletOrg"
                >
                  <i className="bi bi-youtube" />
                </a>
              </li>
            </ul>
          </div>
        </Container>
      </footer>
    </>
  );
}
