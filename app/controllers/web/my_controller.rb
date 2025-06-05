class Web::MyController < Web::ApplicationController
  before_action :authenticate_user!

  def show
    started_language_members = current_user.language_members
      .joins(language: :landing_pages)
      .merge(Language::LandingPage.web.where(main: true))
      .started

    finished_language_members = current_user.language_members
      .joins(language: :landing_pages)
      .merge(Language::LandingPage.web.where(main: true))
      .finished

    landing_pages = Language::LandingPage.web
      .where(main: true)
      .joins(language: :members)
      .merge(current_user.language_members)

    landing_page_resources_by_language_id = landing_pages
      .index_by { it.language_id }
      .transform_values { Language::LandingPageForListsResource.new(it) }

    render inertia: true, props: {
      startedCourseMembers: Language::MemberResource.new(started_language_members),
      finishedCourseMembers: Language::MemberResource.new(finished_language_members),
      landingPageResourcesByCourseId: landing_page_resources_by_language_id
    }
  end
end
