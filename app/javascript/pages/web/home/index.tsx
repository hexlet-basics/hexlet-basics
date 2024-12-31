import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { TypeAnimation } from "react-type-animation";

import Layout from "../../../components/layouts/Application";
import type { Language, LanguageCategory } from "../../../types/serializers";
import * as Routes from "../../../routes.js";

import codeImagePath from "@/images/code-basics-coding-ru.png";
import { assetPath } from "../../../lib/utils.js";

type Props = PropsWithChildren & {
	language_categories: LanguageCategory[];
	languages: Language[];
};

const sequence = [
	"TypeScript",
	1000,
	"Java",
	1000,
	"Python",
	1000,
	"PHP",
	1000,
	"Ruby",
	1000,
	"HTML",
	1000,
	"CSS",
	1000,
	"Go",
	1000,
];

export default function Index({ language_categories, languages }: Props) {
	const { t } = useTranslation();

	return (
		<Layout languageCategories={language_categories} languages={languages}>
			<div className="p-4 pb-0 pt-lg-5 align-items-center border shadow bg-white rounded-3">
				<div className="row">
					<div className="col-lg-7 p-3 p-lg-5 pt-lg-3">
						<h1 className="h6 text-muted">
							{t("home.hero.free_programming_courses")}
						</h1>
						<div className="display-4 fw-bold lh-1 mb-3">
							{t("home.hero.learn_html")}{" "}
							<TypeAnimation
								className="text-primary"
								sequence={sequence}
								wrapper="span"
								speed={5}
								repeat={Number.POSITIVE_INFINITY}
							/>
						</div>
						<p className="lead">{t("home.hero.fastest_way_to_start_coding")}</p>
						<div className="d-grid gap-2 d-md-flex justify-content-md-start mb-4 mb-lg-3">
							<a
								className="btn btn-primary btn-lg px-4 me-md-2 fw-bold"
								href="#courses"
							>
								Попробовать
							</a>
						</div>
					</div>
					<div className="col-lg-4 offset-lg-1 p-0 overflow-hidden shadow-lg text-center">
						<img
							alt="Как работает обучение на code-basics"
							className="rounded-lg-3"
							height="356"
							src={codeImagePath}
							width="720"
						/>
					</div>
				</div>
			</div>

			<div className="my-5">
				<div className="d-flex flex-column flex-sm-row align-items-sm-center mb-4">
					<h2 className="me-auto mb-3 mb-sm-0">
						<a
							id="courses"
							className="text-decoration-none text-dark"
							href="#courses"
						>
							{t("home.languages.courses")}
						</a>
					</h2>

					<div className="d-flex flex-wrap">
						{language_categories.map((c) => (
							<a
								key={c.id}
								href={Routes.language_category_path(c.slug!)}
								className="mb-2 mb-sm-0 me-2 me-sm-0 ms-sm-2 fw-light text-decoration-none fs-5 badge text-bg-light p-2 p-sm-3 rounded-pill border"
							>
								{c.name}
							</a>
						))}
					</div>
				</div>

				<div className="row row-cols-2 row-cols-md-3 row-cols-lg-4">
					{languages.map((language) => (
						<div className="col mb-3" key={language.id}>
							<div className="card h-100 shadow-sm border-0 bg-white">
								<img
									src={assetPath(language.cover!)}
									className="card-img-top"
									alt={language.name!}
								/>
								{/* {languageMembersByLanguage && */}
								{/* 	languageMembersByLanguage[language.id]?.finished && ( */}
								{/* 		<div className="card-img-overlay m-auto bg-light text-center shadow-sm p-1"> */}
								{/* 			{t(".course_finished")} */}
								{/* 		</div> */}
								{/* 	)} */}
								<div className="card-body">
									<h2 className="card-title">
										<a
											href={Routes.language_path(language.slug!)}
											className="stretched-link text-dark text-decoration-none text-nowrap"
										>
											{language.name}
										</a>
									</h2>
									<div className="text-muted">
										<span className="text-nowrap d-inline-block me-4">
											<span className="bi bi-clock me-2" />
											{t("js:hours", { count: language.duration })}
										</span>
										<span className="text-nowrap d-inline-block">
											<span className="bi bi-people me-2" />
											{language.members_count}
										</span>
									</div>
								</div>
							</div>
						</div>
					))}
				</div>
			</div>
		</Layout>
	);
}
