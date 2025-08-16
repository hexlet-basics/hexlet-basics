class Web::CasesController < Web::ApplicationController
  before_action :require_russian_locale

  def index
    seo_tags = {
      title: t(".title"),
      description: t(".meta.description")
    }

    set_meta_tags seo_tags

    render inertia: true, props: {}
  end

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
        site: "@hexlethq"
      }
    }

    set_meta_tags seo_tags

    render inertia: true, props: {}
  end

  # for internal courses
  # for students
  # for programers
end
