# typed: true

class Web::LeadsController < Web::ApplicationController
  before_action :require_authentication

  def new
    seo_tags = {
      title: t(".header")
      # description: t(".meta.description")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      from: params[:from],
      # courseLandingPage: Language::LandingPageResource.new(landing_page)
      lead: LeadCrudResource.new(LeadForm.new)
    }
  end

  def create
    user = T.must(current_user)
    lead = LeadForm.new(params[:data])
    lead.user = user

    courses_data = []
    user.language_members.each do |member|
      course_data = {
        slug: T.must(member.language).slug,
        lessons_finished_count: member.lesson_members.finished.count
      }

      courses_data << course_data
    end

    lead.courses_data = courses_data

    lead.survey_answers_data = []

    lead.ahoy_visit = user.visits.last

    if lead.save

      user.tag_list.remove("should_be_lead")
      user.save!

      lead_created_event_data = {
        lead_id: lead.id,
        user_id: user.id,
        first_name: user.first_name,
        last_name: user.last_name,
        ym_client_id: lead.ym_client_id,
        user_name: user.to_s,
        email: T.must(user.email),
        utm_source: lead.ahoy_visit&.utm_source,
        utm_medium: lead.ahoy_visit&.utm_medium,
        utm_campaign: lead.ahoy_visit&.utm_campaign,
        utm_content: lead.ahoy_visit&.utm_content,
        utm_term: lead.ahoy_visit&.utm_term,
        phone: lead.phone,
        telegram: lead.telegram,
        whatsapp: lead.whatsapp,
        survey_answers_data: lead.survey_answers_data,
        courses_data: lead.courses_data
      }
      lead_created_event = LeadCreatedEvent.new(data: lead_created_event_data)

      publish_event(lead_created_event, current_user)
      js_event(lead_created_event)

      f(:success)
      redirect_to params[:from].presence || root_path
    else
      f(:error)
      redirect_to new_lead_path, inertia: { errors: lead.errors }
    end
  end
end
