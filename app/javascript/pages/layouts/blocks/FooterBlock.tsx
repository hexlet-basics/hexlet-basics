import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";

import { Head, Link, usePage } from "@inertiajs/react";
import { Col, Container, Nav, Row } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import type { Organization, WithContext } from "schema-dts";

export default function FooterBlock() {
  const { t: tLayouts } = useTranslation("layouts");
  const { t: tCommon } = useTranslation("common");
  const { locale, landingPagesForFooter } = usePage<SharedProps>().props;

  const organization: WithContext<Organization> = {
    "@context": "https://schema.org",
    "@type": "Organization",
    email: tCommon("organization.email"),
    name: tCommon("organization.legal_name"),
    telephone: tCommon("organization.phone"),
  };

  return (
    <>
      <Head>
        <script type="application/ld+json">
          {JSON.stringify(organization)}
        </script>
      </Head>
      <footer className="mt-5 bg-body-tertiary pt-5">
        <Container>
          <Row>
            <Col>
              <Nav className="flex-column">
                <Nav.Item className="h5">+7 (495) 085 21 62</Nav.Item>
                <Nav.Item className="h5">8 800 100 22 47</Nav.Item>
                <Nav.Item className="mb-3">
                  <a href={`mailto:${tCommon("organization.email")}`}>
                    {tCommon("organization.email")}
                  </a>
                </Nav.Item>
                <Nav.Item className="mb-3">
                  {tCommon("organization.address")}
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
                {locale === "ru" && (
                  <Nav.Item>
                    <a
                      className="link-body-emphasis text-decoration-none"
                      href="https://ttttt.me/hexlet_ru"
                    >
                      {tLayouts("shared.footer.community")}
                    </a>
                  </Nav.Item>
                )}
                <Nav.Item>
                  <Link
                    className="link-body-emphasis text-decoration-none"
                    href={Routes.map_path()}
                  >
                    {tLayouts("shared.footer.sitemap")}
                  </Link>
                </Nav.Item>
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
              </Nav>
            </Col>
            <Col>
              <Nav className="flex-column">
                {landingPagesForFooter.map((lp) => (
                  <Nav.Item key={lp.id}>
                    <Link
                      className="link-body-emphasis text-decoration-none"
                      href={Routes.language_path(lp.slug)}
                    >
                      {lp.footer_name}
                    </Link>
                  </Nav.Item>
                ))}
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
                <Nav.Item>
                  <Link
                    className="link-body-emphasis text-decoration-none"
                    href={Routes.page_path("cookie_policy")}
                  >
                    {tLayouts("shared.footer.cookie_policy")}
                  </Link>
                </Nav.Item>
              </Nav>
            </Col>
          </Row>
          <div className="d-flex flex-column flex-sm-row justify-content-between pt-4 my-4 border-top">
            <div>{`© ${new Date().getFullYear()} ${tCommon("organization.legal_name")}`}</div>
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
