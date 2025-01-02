import cn from "classnames";
import { type PropsWithChildren, useState } from "react";
import { Col, Container, Nav, NavDropdown, Navbar, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import logoImg from "../../images/logo.png";
import { deviconClass } from "../../lib/utils.js";
import * as Routes from "../../routes.js";

import { Link } from "@inertiajs/react";
import type { Language, LanguageCategory } from "../../types/serializers";

type Props = PropsWithChildren & {
  languageCategories: LanguageCategory[];
  courses: Language[];
};

export default function Application({
  children,
  languageCategories,
  courses,
}: Props) {
  const { t: tLayouts } = useTranslation("layouts");
  const [collapsed, setCollapsed] = useState(false);

  const Logo = <img width="30" alt="Code Basics logo" src={logoImg} />;

  return (
    <>
      <Container className="py-2">
        <Navbar expand="lg" className="mb-5 border-bottom pb-3">
          <Navbar.Brand href="/">
            <img
              src={logoImg}
              width="30"
              height="30"
              className="d-inline-block align-top"
              alt="React Bootstrap logo"
            />
          </Navbar.Brand>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
          <Navbar.Collapse id="basic-navbar-nav">
            <Nav className="me-auto">
              <NavDropdown
                className="link-body-emphasis"
                title={tLayouts("shared.nav.courses")}
                id="basic-nav-dropdown"
              >
                {courses.map((c) => (
                  <NavDropdown.Item
                    className="d-flex align-items-center"
                    key={c.id}
                    href={Routes.language_path(c.slug!)}
                  >
                    <i
                      className={cn(deviconClass(c.slug!), "colored", "me-2")}
                    />
                    {c.name}
                  </NavDropdown.Item>
                ))}
              </NavDropdown>
            </Nav>
            <Nav>
              <Nav.Link
                className="link-body-emphasis"
                href={Routes.new_session_path()}
              >
                {tLayouts("shared.nav.sign_in")}
              </Nav.Link>
              <Nav.Link
                className="link-body-emphasis"
                href={Routes.new_user_path()}
              >
                {tLayouts("shared.nav.registration")}
              </Nav.Link>
            </Nav>
          </Navbar.Collapse>
        </Navbar>
      </Container>
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
                <Nav.Item>
                  ООО «Хекслет Рус», 108813, г. Москва, вн.тер.г. поселение
                  Московский, г. Московский, ул. Солнечная, д. 3А, стр. 1,
                  помещ. 10/3 ОГРН 1217300010476
                </Nav.Item>
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
            <div>© 2025 Хекслет</div>
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
