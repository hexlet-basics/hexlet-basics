import logoImg from "@/images/logo.png";
import defaultAvatarImg from "@/images/user-avatar.png";
import { deviconClass, localesByCode } from "@/lib/utils";
import cn from "classnames";
import i18next from "i18next";
import { Nav, NavDropdown, Navbar, type NavbarProps } from "react-bootstrap";

import useLinkClickHandler from "@/hooks/useLinkClickHandler";
import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";
import type { HTMLAttributes, PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";

type Props = PropsWithChildren & HTMLAttributes<NavbarProps>;

export default function NavbarBlock({ className }: Props) {
  const { courses, auth } = usePage<SharedProps>().props;
  const { t: tLayouts } = useTranslation("layouts");

  const handleLinkClick = useLinkClickHandler();

  return (
    <Navbar expand="lg" className={cn(className, "border-bottom")}>
      <Navbar.Brand href={Routes.root_path()}>
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
          <Nav.Link href={Routes.blog_posts_path()}>{tLayouts('shared.nav.blog')}</Nav.Link>
        </Nav>
        <Nav>
          {auth.user.guest && (
            <>
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
            </>
          )}
          {!auth.user.guest && (
            <NavDropdown
              align="end"
              className="link-body-emphasis"
              title={<i className="bi bi-person-circle" />}
            >
              <NavDropdown.Item href={Routes.edit_account_profile_path()}>
                <div className="d-flex">
                  <img width="50px" src={defaultAvatarImg} alt="User Avatar" />
                  <div>
                    <div className="fw-bold">{auth.user.name}</div>
                    <div>{auth.user.email}</div>
                  </div>
                </div>
              </NavDropdown.Item>
              <NavDropdown.Divider />
              <NavDropdown.Item href={Routes.edit_account_profile_path()}>
                {tLayouts("shared.nav.profile")}
              </NavDropdown.Item>
              {auth.user.admin && (
                <NavDropdown.Item href={Routes.admin_root_path()}>
                  {tLayouts("shared.nav.admin")}
                </NavDropdown.Item>
              )}
              <NavDropdown.Item
                href={Routes.session_path()}
                onClick={handleLinkClick("delete")}
              >
                {tLayouts("shared.nav.sign_out")}
              </NavDropdown.Item>
            </NavDropdown>
          )}
          <NavDropdown
            align="end"
            className="link-body-emphasis"
            title={<i className={localesByCode[i18next.language || "ru"].icon} />}
          >
            {Object.entries(localesByCode).map(([k, v]) => (
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
