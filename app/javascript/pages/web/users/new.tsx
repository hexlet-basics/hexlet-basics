import cn from "classnames";
import { type PropsWithChildren, useState } from "react";
import {
	Button,
	Card,
	Col,
	Container,
	Form,
	FormControl,
	Nav,
	NavDropdown,
	Navbar,
	Row,
} from "react-bootstrap";

import { useTranslation } from "react-i18next";

import logoImg from "../../images/logo.png";
import * as Routes from "@/routes.js";

import { Link } from "@inertiajs/react";
import type { Language, LanguageCategory, User } from "@/types/serializers";
import Application from "@/pages/layouts/Application";
import { XInput } from "@/components/forms";

type Props = PropsWithChildren & {
	languageCategories: LanguageCategory[];
	languages: Language[];
	user: User;
};

export default function New({ languageCategories, languages, user }: Props) {
	const { t } = useTranslation();
	const { t: tHelpers } = useTranslation("helpers");
	const { t: tAr } = useTranslation("activerecord");

	return (
		<Application languageCategories={languageCategories} languages={languages}>
			<Container>
				<Row className="justify-content-center">
					<div className="col-sm-8 col-md-7 col-lg-5">
						<h1 className="text-center mb-3">{t("users.new.sign_up")}</h1>
						<Card className="p-4 border-0">
							<Card.Body>
								<Form className="d-flex flex-column">
									<XInput model={user} attribute="first_name" />
									<XInput model={user} attribute="email" />
									<div className="text-end text-muted small mb-4">
										{t("users.new.have_account")}{" "}
										<Link
											href={Routes.new_session_path()}
											className="text-decoration-none"
										>
											{t("users.new.sign_in")}
										</Link>
									</div>
									<Button className="mb-3" type="submit">
										{tHelpers("submit.user_sign_up_form.create")}
									</Button>
									<div
										className="small text-muted"
										dangerouslySetInnerHTML={{
											__html: t("users.new.confirmation_html"),
										}}
									/>
								</Form>
							</Card.Body>
						</Card>
					</div>
				</Row>
			</Container>
		</Application>
	);
}
