import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";

import Application from "@/pages/layouts/Application";
import type {
	BlogPost,
	Language,
	LanguageCategory,
	User,
} from "@/types/serializers";
import i18next from "i18next";
import Markdown from "react-markdown";
import { XBreadcrumb } from "@/components/breadcrumbs";
import { BreadcrumbItem } from "@/types/types";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

type Props = PropsWithChildren & {
	languageCategories: LanguageCategory[];
	languages: Language[];
	blogPost: BlogPost;
	recommendedBlogPosts: BlogPost[];
};

export default function New({
	languageCategories,
	languages,
	blogPost,
	recommendedBlogPosts,
}: Props) {
	const { t } = useTranslation();

	const items: BreadcrumbItem[] = [
		{
			name: t("blog_posts.index.header"),
			url: Routes.blog_posts_path(),
		},
		{
			name: blogPost.name!,
			url: Routes.blog_post_path(blogPost.slug!),
		},
	];

	return (
		<Application languageCategories={languageCategories} languages={languages}>
			<Container>
				<Row className="justify-content-center">
					<Col className="col-12 col-md-10 col-lg-7">
						<XBreadcrumb items={items} />
						<h1>{blogPost.name}</h1>
						<Markdown>{blogPost.body}</Markdown>
					</Col>
				</Row>
			</Container>
		</Application>
	);
}
