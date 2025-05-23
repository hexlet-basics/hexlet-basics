class Web::ApplicationController < ApplicationController
  include BrowserConcern
  include SurveyConcern

  inertia_share flash: -> { flash.to_hash }

  include FlashConcern
  include EventConcern
  include LocaleConcern

  before_action :prepare_locale_settings
  before_action :run_survey_if_needed

  before_action do
    set_meta_tags site: "CodeBasics"
  end

  inertia_share do
    language_categories = Language::Category.with_locale
    landing_pages_for_lists = Language::LandingPage.web
      .where(main: true).where(listed: true)
      .merge(Language.ordered)

    landing_pages_for_footer = Language::LandingPage.web
      .where(footer: true)
      .merge(Language.ordered)

    {
      shouldAddContactMethod: current_user.should_be_lead,
      courseCategories: Language::CategoryResource.new(language_categories),
      railsDirectUploadsUrl: view_context.rails_direct_uploads_url,
      landingPagesForLists: Language::LandingPageForListsResource.new(landing_pages_for_lists),
      landingPagesForFooter: Language::LandingPageForListsResource.new(landing_pages_for_footer),
      locale: I18n.locale,
      suffix: I18n.locale == :en ? nil : I18n.locale,
      auth: {
        user: UserResource.new(current_user)
      },
      mobileBrowser: mobile_browser?,
      # carrotQuestUserHash: signed_in? ? OpenSSL::HMAC.hexdigest("SHA256", configus.carrotquest_user_auth_key, current_user.id.to_s) : nil,
      metaTagsHTMLString: display_escaped_meta_tags(reverse: true)
    }
  end

  before_action do
    gon.current_user = UserResource.new(current_user)
    gon.locale = I18n.locale
    gon.suffix = I18n.locale == :en ? nil : I18n.locale
  end

  private

  # def check!(value)
  #   raise ActiveRecord::RecordNotFound if value.nil?
  # end

  def ransack_params(defaults)
    raw = params.permit(:sf, :so, fields: {}).with_defaults({ fields: {} }).with_defaults(defaults)
    ransack = raw["fields"]

    if raw.key?("sf")
      ransack["s"] = "#{raw["sf"]} #{raw["so"] == "1" ? 'asc' : 'desc'}"
    end

    ransack.to_unsafe_hash
  end

  def grid_params(pagy = nil)
    result = params.permit(:sf, :so, fields: {})
    if pagy
      result[:page] = pagy.page
      result[:tr] = pagy.count()
      result[:per] = pagy.limit()
    else
      result[:page] = 1
      result[:tr] = 0
      result[:per] = 10
    end

    OpenStruct.new(result)
  end

  def redirect_to_inertia(url, model)
    # if model.errors.any? && Rails.env.test?
    #   raise model.errors.full_messages.inspect
    # end
    redirect_to url, inertia: { errors: model.errors }
  end

  def escape_meta_tags(tags)
    tags.transform_values! do |tag|
      if tag.is_a?(Hash)
        escape_meta_tags(tag)
      elsif tag.is_a?(Array)
        tag.map! { |item| Rack::Utils.escape_html(item) }
      else
        Rack::Utils.escape_html(tag)
      end
    end
  end

  def display_escaped_meta_tags(...)
    escape_meta_tags(meta_tags.instance_values["meta_tags"])
    helpers.display_meta_tags(...)
  end
end
