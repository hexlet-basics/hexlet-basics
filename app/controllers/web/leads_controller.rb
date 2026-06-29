# typed: strict

class Web::LeadsController < Web::ApplicationController
  before_action :require_authentication

  sig { returns(T.untyped) }
  def new
    seo_tags = {
      title: t(".header")
      # description: t(".meta.description")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      from: params[:from],
      # courseLandingPage: Language::LandingPageResource.new(landing_page)
      lead: LeadCrudResource.new(LeadStruct.new)
    }
  end

  sig { returns(T.untyped) }
  def create
    user = T.must(current_user)
    struct = ApplicationParamsStruct.from_params(LeadStruct, params.require(:data))

    return fail_create(struct.errors) unless struct.valid?

    courses_data = user.language_members.map do |member|
      {
        slug: member.language.slug,
        lessons_finished_count: member.lesson_members.finished.count
      }
    end

    result = LeadService.create(struct, user:, courses_data:, ahoy_visit: user.visits.last)

    case result
    when Typed::Failure
      return fail_create(result.error.errors)
    end

    lead = result.payload

    user.tag_list.remove("should_be_lead")
    user.save!

    lead_created_event_data = {
      lead_id: lead.id,
      user_id: user.id,
      first_name: user.first_name,
      last_name: user.last_name,
      ym_client_id: lead.ym_client_id,
      user_name: user.to_s,
      email: user.email,
      utm_source: lead.ahoy_visit&.utm_source,
      utm_medium: lead.ahoy_visit&.utm_medium,
      utm_campaign: lead.ahoy_visit&.utm_campaign,
      utm_content: lead.ahoy_visit&.utm_content,
      utm_term: lead.ahoy_visit&.utm_term,
      phone: lead.phone,
      telegram: lead.telegram,
      whatsapp: lead.whatsapp
    }
    lead_created_event = LeadCreatedEvent.new(data: lead_created_event_data)

    publish_event(lead_created_event, user)
    js_event(lead_created_event)

    f(:success)
    redirect_to params[:from].presence || root_path
  end

  private

  sig { params(errors: T.untyped).returns(T.untyped) }
  def fail_create(errors)
    f(:error)
    redirect_to new_lead_path, inertia: { errors: }
  end
end
