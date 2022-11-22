# frozen_string_literal: true

class Web::HomeController < Web::ApplicationController
  def index
    @languages_completed = Language.with_progress(:completed)
                                   .joins(current_version: :infos)
                                   .includes(:current_version)
                                   .merge(Language::Version::Info.with_locale)
                                   .ordered

    @languages_in_development = Language.with_progress(:in_development).includes(:current_version)
    @language_members_by_language = current_user.language_members.index_by(&:language_id)

    language_versions = Language::Version.where(current_language: @languages_completed)

    @js_course = Language.find_by slug: 'javascript'
    @html_course = Language.find_by slug: 'html'

    @user = User.new
    @users_count = User.count

    @languages_links_by_slug = Language.all.each_with_object({}) do |item, acc|
      acc[item.slug.to_sym] = view_context.link_to(item, language_path(item.slug))
    end

    @categories = Language::Category.all

    gon.languages_for_widget = language_versions.pluck(:name)

    seo_tags = {
      title: t('.title'),
      description: t('.meta_description'),
      canonical: root_url,
      image_src: view_context.asset_url('logo.png'),
      twitter: {
        card: 'summary',
        site: '@hexlet_io'
      },
      og: {
        title: t('.title'),
        type: 'website',
        url: root_url,
        image: view_context.asset_url('logo.png')
      }
    }
    set_meta_tags seo_tags
  end

  def robots
    respond_to :text
  end
end
