# typed: strict

class Web::CasesController < Web::ApplicationController
  allow_unauthenticated_access
  before_action :require_russian_locale

  sig { void }
  def index
    seo_tags = {
      title: t(".title"),
      description: t(".meta.description")
    }

    set_meta_tags seo_tags

    render inertia: true, props: {}
  end

  sig { void }
  def for_teachers
    description = t(".meta.description").truncate(160)

    seo_tags = {
      title: t(".title"),
      description: description,
      canonical: view_context.book_url,
      og: {
        title: t(".title"),
        description:  description
      },
      twitter: {
        card: "summary",
        site: t("links.hexlet_twitter_handle")
      }
    }

    set_meta_tags seo_tags

    render inertia: true, props: {}
  end

  # for internal courses
  # for students
  # for programers
end
