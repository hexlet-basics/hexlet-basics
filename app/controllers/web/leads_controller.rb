class Web::LeadsController < Web::ApplicationController
  before_action :authenticate_user!

  def new
    seo_tags = {
      title: t(".header")
      # description: t(".meta.description")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      # courseLandingPage: Language::LandingPageResource.new(landing_page)
      lead: LeadCrudResource.new(LeadForm.new)
    }
  end

  def create
    lead = LeadForm.new(params[:lead])
    lead.user = current_user

    courses_data = []
    current_user.language_members.each do |member|
      course_data = {
        slug: member.language.slug,
        lessons_finished_count: member.lesson_members.finished.count
      }

      courses_data << course_data
    end

    lead.courses_data = courses_data

    survey_answers_data = []
    current_user.survey_answers.each do |answer|
      answer_data = {
        question: answer.survey.question,
        answer: answer.survey_item.value
      }

      survey_answers_data << answer_data
    end

    lead.survey_answers_data = survey_answers_data

    if lead.save

      current_user.tag_list.remove("should_be_lead")
      current_user.save!

      lead_created_event_data = {
        user_id: current_user.id,
        first_name: current_user.first_name,
        last_name: current_user.last_name,
        ym_client_id: lead.ym_client_id,
        user_name: current_user.to_s,
        email: current_user.email,
        phone: lead.phone,
        telegram: lead.telegram,
        whatsapp: lead.whatsapp,
        survey_answers_data: lead.survey_answers_data,
        courses_data: lead.courses_data
      }
      lead_created_event = LeadCreatedEvent.new(data: lead_created_event_data)

      publish_event(lead_created_event, current_user)
      event_to_js(lead_created_event)

      f(:success)
      redirect_to params[:from].presence || root_path
    else
      f(:error)
      redirect_to_inertia new_lead_path, lead
    end
  end
end
