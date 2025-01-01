import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import Markdown from "react-markdown";
import { TypeAnimation } from "react-type-animation";

import Layout from "@/components/layouts/Application";
import * as Routes from "@/routes.js";
import type {
	BlogPost,
	Language,
	LanguageCategory,
	Review,
	User,
} from "@/types/serializers";

import codeImagePath from "@/images/code-basics-coding-ru.png";
import { Accordion, Container } from "react-bootstrap";
import { assetPath } from "@/lib/utils.js";

type Props = PropsWithChildren & {
	languageCategories: LanguageCategory[];
	languages: Language[];
	blogPosts: BlogPost[];
	reviews: Review[];
	currentUser: User;
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

const reviews = [
	{
		name: "Александр Авдошкин",
		avatar: assetPath("avdoshkin.jpg"),
		body: `Если бы не коронавирус, выполнил бы всё в заход (в смысле каждый день по несколько пунктов в теме).
            Изучаю с нуля, ваш портал очень ориентирован на новичков. Спасибо вам большое!`,
	},
	{
		name: "Сергей Тюрин",
		avatar: assetPath("tyrin.jpg"),
		body: `Очень всё доступно даже для полного профана вроде меня. Эта вводная по JS вошла в мой туговатый ум,
            складно как недостающий пазл. Всем кидаю линк на эту страничку.`,
	},
	{
		name: "Элиях Клейман",
		avatar: assetPath("user-avatar.png"),
		body: `Для меня это первый курс для новичка. Понравилось тем, что вся информация структурирована
            и дана по мере изучения материала в иерархичном порядке, что значительно повышает и желание к обучению`,
	},
];

export default function Index({
	languageCategories,
	languages,
	blogPosts,
	currentUser,
}: Props) {
	const { t } = useTranslation();
	const { t: tFaq } = useTranslation("faq");
	const { t: tJS } = useTranslation("js");

	const faq = tFaq("main", { returnObjects: true });
	console.log(faq);

	return (
		<Layout languageCategories={languageCategories} languages={languages}>
			<Container className="mb-5">
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
							<p className="lead">
								{t("home.hero.fastest_way_to_start_coding")}
							</p>
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
			</Container>

			<Container className="mb-5">
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
						{languageCategories.map((c) => (
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
								{/*   languageMembersByLanguage[language.id]?.finished && ( */}
								{/*     <div className="card-img-overlay m-auto bg-light text-center shadow-sm p-1"> */}
								{/*       {t(".course_finished")} */}
								{/*     </div> */}
								{/*   )} */}
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
											{tJS("hours", { count: language.duration })}
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
			</Container>

			<Container className="mb-5">
				<div className="d-flex">
					<h2 className="me-auto mt-auto">{t("home.index.reviews")}</h2>
					<div className="mt-auto">
						<a
							href={Routes.reviews_path()}
							className="text-decoration-none text-muted small"
						>
							{t("home.index.all_reviews")}
						</a>
					</div>
				</div>
				<hr className="mb-5 mt-1" />

				<div className="row g-4 row-cols-1 row-cols-sm-2 row-cols-lg-3">
					{reviews.map((review) => (
						<div key={review.avatar} className="col">
							<div className="d-flex mb-3">
								<img
									src={review.avatar}
									className="rounded-circle flex-shrink-0"
									width={50}
									height={50}
									alt={`Аватар пользователя ${review.name}`}
								/>
								<div className="ms-3">
									<div className="fw-bold">{review.name}</div>
									<div className="small">
										{/* {t(".course_html", { */}
										{/*   link: ( */}
										{/*     <a */}
										{/*       href={language_path(jsCourse.slug)} */}
										{/*       className="text-dark" */}
										{/*     > */}
										{/*       {jsCourse} */}
										{/*     </a> */}
										{/*   ), */}
										{/* })} */}
									</div>
								</div>
							</div>
							<div className="fst-italic">{review.body}</div>
						</div>
					))}
				</div>
			</Container>

			<Container className="mb-5">
				{blogPosts.length > 0 && (
					<>
						<div className="d-flex">
							<h2 className="me-auto mt-auto">
								{t("home.index.blog_posts")}
							</h2>
							<div className="mt-auto">
								<a
									href={Routes.blog_posts_path()}
									className="text-decoration-none text-muted small"
								>
									{t("home.index.all_blog_posts")}
								</a>
							</div>
						</div>
						<hr className="mb-5 mt-1" />
						<div className="row row-cols-sm-2 row-cols-md-3 row-cols-1">
							{blogPosts.map((post) => (
								<div className="col" key={post.id}>
									<div className="card border-0 bg-body">
										{post.cover && (
											<a href={Routes.blog_post_path(post.slug!)}>
												{/* <img */}
												{/*   src={post.cover.variant("list")} */}
												{/*   className="card-img-top img-fluid" */}
												{/*   alt={`Cover for ${post.title}`} */}
												{/* /> */}
											</a>
										)}
										<div className="card-body">
											<h2 className="card-title fs-6">
												<a
													href={Routes.blog_post_path(post.slug!)}
													className="text-decoration-none"
												>
													{post.name}
												</a>
											</h2>
											<div className="card-text">{post.description}</div>
										</div>
									</div>
								</div>
							))}
						</div>
					</>
				)}
			</Container>

			<Container className="mb-5 pb-5">
				<h2>
					<a id="faq" className="text-decoration-none text-dark" href="#faq">
						{tFaq("header")}
					</a>
				</h2>
				<hr className="mb-5" />
				<Accordion defaultActiveKey="0">
					{Object.entries(faq).map(([key, value], index) => (
						<Accordion.Item eventKey={String(index)} key={key}>
							<Accordion.Header>{value.question}</Accordion.Header>
							<Accordion.Body>
								<Markdown>{value.answer}</Markdown>
							</Accordion.Body>
						</Accordion.Item>
					))}
				</Accordion>
			</Container>

			{currentUser.guest && (
				<div className="bg-white">
					<div className="container">
						<div className="row align-items-center g-lg-5 py-5">
							<div className="col-lg-7 text-center text-lg-start">
								{/* {cache([I18n.locale], { expiresIn: "1d" }, () => ( */}
								{/*   <> */}
								{/*     <h1 className="display-4 fw-bold lh-1 mb-3"> */}
								{/*       {t(".join", { count: User.count })} */}
								{/*     </h1> */}
								{/*     Uncomment this if the paragraph is needed */}
								{/*     <p className="col-lg-10 fs-4">Forever free.</p> */}
								{/*   </> */}
								{/* ))} */}
							</div>
							<div className="col-md-10 mx-auto col-lg-5">
								<div className="card p-4">
									<div className="card-body">
										<form
											action={Routes.users_path()}
											method="POST"
											className="simple_form"
										>
											<div className="form-group">
												<label htmlFor="user_first_name">First Name</label>
												<input
													type="text"
													id="user_first_name"
													name="user[first_name]"
													className="form-control"
												/>
											</div>
											<div className="form-group">
												<label htmlFor="user_email">Email</label>
												<input
													type="email"
													id="user_email"
													name="user[email]"
													autoComplete="email"
													className="form-control"
												/>
											</div>
											<div className="form-group">
												<label htmlFor="user_password">Password</label>
												<input
													type="password"
													id="user_password"
													name="user[password]"
													autoComplete="new-password"
													className="form-control"
												/>
											</div>
											<div className="form-group">
												<button type="submit" className="btn btn-primary w-100">
													Submit
												</button>
											</div>
										</form>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			)}
		</Layout>
	);
}
