import cn from "classnames";
import { type PropsWithChildren, useState } from "react";
import { Container, Nav, Navbar, NavDropdown } from "react-bootstrap";

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
		<Container>
			<Navbar expand="lg" className="mb-5 border-bottom">
				<Container>
					<Navbar.Brand href="#home">
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
									<Nav.Link key={c.id} href={Routes.language_path(c.slug!)}>
										{c.name}
									</Nav.Link>
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
      {children}
		</Container>
		// <main className="container mx-auto px-4">
		//   {/* <div className="navbar bg-base-100"> */}
		//   {/*   <div className="navbar-start"> */}
		//   {/*     <a href="/" className="btn btn-ghost text-xl"> */}
		//   {/*       CodeBasics */}
		//   {/*     </a> */}
		//   {/*     <ul className="menu menu-horizontal"> */}
		//   {/*       <li> */}
		//   {/*         <details> */}
		//   {/*           <summary>{t("shared.nav.courses")}</summary> */}
		//   {/*           <ul className="bg-base-100 rounded-t-none p-2"> */}
		//   {/*             {languages.map((c) => ( */}
		//   {/*               <li key={c.id}> */}
		//   {/*                 <a href={Routes.language_path(c.slug!)}>{c.name}</a> */}
		//   {/*               </li> */}
		//   {/*             ))} */}
		//   {/*           </ul> */}
		//   {/*         </details> */}
		//   {/*       </li> */}
		//   {/*       <li> */}
		//   {/*         <details> */}
		//   {/*           <summary>{t("shared.nav.categories")}</summary> */}
		//   {/*           <ul className="bg-base-100 rounded-t-none p-2"> */}
		//   {/*             {languageCategories.map((c) => ( */}
		//   {/*               <li key={c.id}> */}
		//   {/*                 <a href={Routes.language_category_path(c.slug!)}> */}
		//   {/*                   {c.name} */}
		//   {/*                 </a> */}
		//   {/*               </li> */}
		//   {/*             ))} */}
		//   {/*           </ul> */}
		//   {/*         </details> */}
		//   {/*       </li> */}
		//   {/*     </ul> */}
		//   {/*   </div> */}
		//   {/*   <div className="navbar-end"> */}
		//   {/*     <ul className="menu menu-horizontal"> */}
		//   {/*       <li> */}
		//   {/*         <a href={Routes.new_session_path()}>{t("shared.nav.sign_in")}</a> */}
		//   {/*       </li> */}
		//   {/*       <li> */}
		//   {/*         <a href={Routes.new_user_path()}> */}
		//   {/*           {t("shared.nav.registration")} */}
		//   {/*         </a> */}
		//   {/*       </li> */}
		//   {/*     </ul> */}
		//   {/*   </div> */}
		//   {/* </div> */}
		//   {/* <header className="container mb-4 flex-shrink-0"> */}
		//   {/*   <div className="py-2 border-bottom"> */}
		//   {/*     <div className="navbar navbar-expand-md navbar-light"> */}
		//   {/*       <div className="container-fluid"> */}
		//   {/*         <a className="navbar-brand me-md-4" href="/ru"> */}
		//   {/*           <img */}
		//   {/*             width="30" */}
		//   {/*             alt="Code Basics logo" */}
		//   {/*             src={logoImg} */}
		//   {/*             loading="lazy" */}
		//   {/*           /> */}
		//   {/*         </a> */}
		//   {/*         <button */}
		//   {/*           aria-controls="navbarToggler" */}
		//   {/*           aria-expanded="false" */}
		//   {/*           aria-label="Переключить навигацию" */}
		//   {/*           className="navbar-toggler" */}
		//   {/*           data-bs-target="#navbarToggler" */}
		//   {/*           data-bs-toggle="collapse" */}
		//   {/*           type="button" */}
		//   {/*         > */}
		//   {/*           <span className="navbar-toggler-icon" /> */}
		//   {/*         </button> */}
		//   {/*         <div className="collapse navbar-collapse" id="navbarToggler"> */}
		//   {/*           <ul className="navbar-nav me-auto"> */}
		//   {/*             <li className="nav-item dropdown me-2"> */}
		//   {/*               <a */}
		//   {/*                 aria-expanded="false" */}
		//   {/*                 aria-haspopup="true" */}
		//   {/*                 className="nav-link text-dark dropdown-toggle" */}
		//   {/*                 data-bs-toggle="dropdown" */}
		//   {/*                 href="#" */}
		//   {/*               > */}
		//   {/*                 Курсы */}
		//   {/*               </a> */}
		//   {/*               <div className="dropdown-menu x-mw-320"></div> */}
		//   {/*             </li> */}
		//   {/*             <li className="nav-item dropdown me-2"> */}
		//   {/*               <a */}
		//   {/*                 aria-expanded="false" */}
		//   {/*                 aria-haspopup="true" */}
		//   {/*                 className="nav-link text-dark dropdown-toggle" */}
		//   {/*                 data-bs-toggle="dropdown" */}
		//   {/*                 href="#" */}
		//   {/*               > */}
		//   {/*                 Категории */}
		//   {/*               </a> */}
		//   {/*               <ul className="dropdown-menu dropdown-menu-columns dropdown-menu-arrow"> */}
		//   {/*                 <li> */}
		//   {/*                   <a */}
		//   {/*                     className="dropdown-item" */}
		//   {/*                     href="/ru/language_categories/programming" */}
		//   {/*                   > */}
		//   {/*                     <span>Программирование</span> */}
		//   {/*                   </a> */}
		//   {/*                 </li> */}
		//   {/*                 <li> */}
		//   {/*                   <a */}
		//   {/*                     className="dropdown-item" */}
		//   {/*                     href="/ru/language_categories/layouting" */}
		//   {/*                   > */}
		//   {/*                     <span>Верстка</span> */}
		//   {/*                   </a> */}
		//   {/*                 </li> */}
		//   {/*               </ul> */}
		//   {/*             </li> */}
		//   {/*             <li className="nav-item dropdown me-2"> */}
		//   {/*               <a */}
		//   {/*                 aria-expanded="false" */}
		//   {/*                 aria-haspopup="true" */}
		//   {/*                 className="nav-link text-dark dropdown-toggle" */}
		//   {/*                 data-bs-toggle="dropdown" */}
		//   {/*                 href="#" */}
		//   {/*               > */}
		//   {/*                 О Code Basics */}
		//   {/*               </a> */}
		//   {/*               <ul className="dropdown-menu"> */}
		//   {/*                 <li> */}
		//   {/*                   <a */}
		//   {/*                     className="dropdown-item" */}
		//   {/*                     href={Routes.reviews_path()} */}
		//   {/*                   > */}
		//   {/*                     Отзывы */}
		//   {/*                   </a> */}
		//   {/*                 </li> */}
		//   {/*                 <li> */}
		//   {/*                   <a */}
		//   {/*                     className="dropdown-item" */}
		//   {/*                     href={Routes.page_path("about")} */}
		//   {/*                   > */}
		//   {/*                     О платформе */}
		//   {/*                   </a> */}
		//   {/*                 </li> */}
		//   {/*                 <li> */}
		//   {/*                   <a */}
		//   {/*                     className="dropdown-item" */}
		//   {/*                     href={Routes.blog_posts_path()} */}
		//   {/*                   > */}
		//   {/*                     Блог */}
		//   {/*                   </a> */}
		//   {/*                 </li> */}
		//   {/*                 <li> */}
		//   {/*                   <a */}
		//   {/*                     className="dropdown-item" */}
		//   {/*                     target="_blank" */}
		//   {/*                     rel="nofollow noopener noreferrer" */}
		//   {/*                     href="https://ru.hexlet.io/blog/categories/success?promo_name=blog-success&amp;promo_position=body&amp;promo_type=link" */}
		//   {/*                   > */}
		//   {/*                     Истории успеха */}
		//   {/*                   </a> */}
		//   {/*                 </li> */}
		//   {/*                 <li> */}
		//   {/*                   <a className="dropdown-item" href="/ru/pages/authors"> */}
		//   {/*                     Авторам */}
		//   {/*                   </a> */}
		//   {/*                 </li> */}
		//   {/*               </ul> */}
		//   {/*             </li> */}
		//   {/*           </ul> */}
		//   {/*           <ul className="navbar-nav"> */}
		//   {/*             <li className="nav-item"> */}
		//   {/*               <a className="nav-link  " href="/ru/session/new"> */}
		//   {/*                 Вход */}
		//   {/*               </a> */}
		//   {/*             </li> */}
		//   {/*             <li className="nav-item"> */}
		//   {/*               <a className="nav-link me-3 " href="/ru/users/new"> */}
		//   {/*                 Регистрация */}
		//   {/*               </a> */}
		//   {/*             </li> */}
		//   {/*             <li className="nav-item dropdown"> */}
		//   {/*               <a */}
		//   {/*                 aria-expanded="false" */}
		//   {/*                 aria-haspopup="true" */}
		//   {/*                 className="nav-link text-dark dropdown-toggle" */}
		//   {/*                 data-bs-toggle="dropdown" */}
		//   {/*                 href="#" */}
		//   {/*               > */}
		//   {/*                 <span className="me-2">RU</span> */}
		//   {/*               </a> */}
		//   {/*               <ul className="dropdown-menu dropdown-menu-end"> */}
		//   {/*                 <li> */}
		//   {/*                   <a */}
		//   {/*                     className="dropdown-item d-flex" */}
		//   {/*                     rel="nofollow" */}
		//   {/*                     href="/ru/locale/switch?new_locale=en" */}
		//   {/*                   > */}
		//   {/*                     <img */}
		//   {/*                       width="25" */}
		//   {/*                       alt="switch language to en" */}
		//   {/*                       className="my-auto me-2" */}
		//   {/*                       src="" */}
		//   {/*                       loading="lazy" */}
		//   {/*                     /> */}
		//   {/*                     <span className="my-auto">Английский</span> */}
		//   {/*                   </a> */}
		//   {/*                 </li> */}
		//   {/*               </ul> */}
		//   {/*             </li> */}
		//   {/*           </ul> */}
		//   {/*         </div> */}
		//   {/*       </div> */}
		//   {/*     </div> */}
		//   {/*   </div> */}
		//   {/* </header> */}
		//   <div>{children}</div>
		// </main>
	);
}
