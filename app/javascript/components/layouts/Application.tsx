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
	languages: Language[];
};

export default function Application({
	children,
	languageCategories,
	languages,
}: Props) {
	const { t } = useTranslation("layouts");
	const [collapsed, setCollapsed] = useState(false);

	const Logo = <img width="30" alt="Code Basics logo" src={logoImg} />;

	return (
		<>
			<Container>
				<Navbar expand="lg" className="mb-5 border-bottom">
					<Container>
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
									title={t("shared.nav.courses")}
									id="basic-nav-dropdown"
								>
									{languages.map((c) => (
										<NavDropdown.Item className="d-flex align-items-center" key={c.id} href={Routes.language_path(c.slug!)}>
                      <i className={cn(deviconClass(c.slug!), 'colored', 'me-2')} />
											{c.name}
										</NavDropdown.Item>
									))}
								</NavDropdown>
							</Nav>
							<Nav>
								<Nav.Link href={Routes.new_session_path()}>
									{t("shared.nav.sign_in")}
								</Nav.Link>
								<Nav.Link href={Routes.new_user_path()}>
									{t("shared.nav.registration")}
								</Nav.Link>
							</Nav>
						</Navbar.Collapse>
					</Container>
				</Navbar>
			</Container>
			{children}
			<footer className="mt-5">
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
							<Nav>
								<Nav.Item>
									<Link href={""}>Сообщество</Link>
								</Nav.Item>
							</Nav>
						</Col>
						<Col>
							<Nav>
								<Nav.Item>
									<Link href={""}>Сообщество</Link>
								</Nav.Item>
							</Nav>
						</Col>
					</Row>
					<div className="d-flex flex-column flex-sm-row justify-content-between pt-4 my-4 border-top">
						<div>© 2025 Хекслет</div>
						<ul className="fs-3 d-flex list-unstyled">
							<li className="me-3">
								<a target="_blank" rel="noreferrer" className="link-body-emphasis" href="https://ttttt.me/hexlet_ru">
									<i className="bi bi-telegram" />
								</a>
							</li>
							<li>
								<a target="_blank" rel="noreferrer" className="link-body-emphasis" href="https://www.youtube.com/@HexletOrg">
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
