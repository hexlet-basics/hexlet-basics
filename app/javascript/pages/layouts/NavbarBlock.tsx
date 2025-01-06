import i18next from "i18next";
import cn from "classnames";
import logoImg from "@/images/logo.png";
import { deviconClass } from "@/lib/utils";
import { Nav, NavDropdown, Navbar, type NavbarProps } from "react-bootstrap";

import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types/types.ts";
import { usePage } from "@inertiajs/react";
import { useTranslation } from "react-i18next";
import type { HTMLAttributes, PropsWithChildren } from "react";

const locales = {
  ru: {
    icon: "fi fi-ru",
    name: "Русский",
  },
  en: {
    icon: "fi fi-us",
    name: "English",
  },
};

type Props = PropsWithChildren & HTMLAttributes<NavbarProps>;

export default function NavbarBlock({ className }: Props) {
  const { courses } = usePage<SharedProps>().props;
  const { t: tLayouts } = useTranslation("layouts");

  return (
    <Navbar expand="lg" className={cn(className, "border-bottom")}>
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
                <i className={cn(deviconClass(c.slug!), "colored", "me-2")} />
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
          <NavDropdown
            drop="start"
            className="link-body-emphasis"
            title={<i className={locales[i18next.language].icon} />}
            id="basic-nav-dropdown"
          >
            {Object.entries(locales).map(([k, v]) => (
              <NavDropdown.Item
                className="d-flex align-items-center"
                key={k}
                href={Routes.switch_locale_path({ new_locale: k })}
              >
                <i className={cn(v.icon, "me-2")} />
                {v.name}
              </NavDropdown.Item>
            ))}
          </NavDropdown>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  );
}
