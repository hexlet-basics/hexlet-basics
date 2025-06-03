class Web::CasesController < Web::ApplicationController
  def index
    seo_tags = {
      title: t(".title"),
      description: t(".meta.description")
    }

    set_meta_tags seo_tags

    render inertia: true, props: {}
  end

  def for_teachers
    seo_tags = {
      title: t(".title"),
      description: t(".meta.description")
    }

    set_meta_tags seo_tags

    render inertia: true, props: {}
  end

  # for internal courses
  # for students
  # for programers
end
