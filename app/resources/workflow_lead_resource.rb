# typed: strict

class WorkflowLeadResource
  extend T::Sig

  LEAD_FIELD_MAPPER = {
    ym_client_id: 316_941, # _ym_uid
    utm_source: 316_919,
    utm_medium: 316_915,
    utm_campaign: 316_917,
    utm_content: 316_913,
    utm_term: 316_921
  }.freeze

  CONTACT_FIELD_MAPPER = {
    user_id: 1_030_881,
    email: { field_id: 316_907, enum_id: 158_463 }, # Личная почта
    phone: { field_id: 316_905, enum_id: 158_453 }, # Мобильный телефон
    telegram: 925_081,
    whatsapp: 1_012_781
  }.freeze

  sig { params(lead: Lead).void }
  def initialize(lead)
    @lead = lead
  end

  sig { returns(T::Hash[String, T.untyped]) }
  def to_h
    {
      "lead" => lead_payload,
      "note" => { "text" => note_text }
    }
  end

  private

  sig { returns(T::Array[T::Hash[Symbol, T.untyped]]) }
  def lead_payload
    lead_entity = {
      name: @lead.email.presence || @lead.user.email,
      pipeline_id: 9_614_774,
      status_id: 76_748_378,
      responsible_user_id: 7_877_026
    }

    custom_fields_values = build_custom_fields(lead_data, LEAD_FIELD_MAPPER)
    lead_entity[:custom_fields_values] = custom_fields_values unless custom_fields_values.empty?

    [
      {
        source_uid: "lead_form",
        source_name: "lead_form",
        metadata: metadata_payload,
        _embedded: {
          contacts: [ contact_payload ],
          leads: [ lead_entity ]
        }
      }
    ]
  end

  sig { returns(T::Hash[Symbol, T.untyped]) }
  def metadata_payload
    visit = @lead.ahoy_visit
    ip = visit&.ip.to_s
    is_ipv4 = ip.present? && IPAddr.new(ip).ipv4?

    {
      form_id: "lead_form",
      form_name: "lead_form",
      form_sent_at: @lead.created_at.to_i,
      form_page: visit&.landing_page,
      ip: is_ipv4 ? ip : nil,
      referer: visit&.referrer,
      visitor_uid: visit&.visitor_token
    }.compact
  end

  sig { returns(T::Hash[Symbol, T.untyped]) }
  def contact_payload
    user = @lead.user
    data = {
      user_id: user.id,
      email: @lead.email.presence || user.email,
      phone: @lead.phone,
      telegram: @lead.telegram,
      whatsapp: @lead.whatsapp
    }

    {
      name: user.to_s,
      first_name: user.first_name.to_s,
      last_name: user.last_name.to_s,
      custom_fields_values: build_custom_fields(data, CONTACT_FIELD_MAPPER)
    }
  end

  sig { returns(T::Hash[Symbol, T.untyped]) }
  def lead_data
    visit = @lead.ahoy_visit

    {
      ym_client_id: @lead.ym_client_id.presence || visit&.ym_client_id,
      utm_source: visit&.utm_source,
      utm_medium: visit&.utm_medium,
      utm_campaign: visit&.utm_campaign,
      utm_content: visit&.utm_content,
      utm_term: visit&.utm_term
    }
  end

  sig do
    params(
      data: T::Hash[Symbol, T.untyped],
      mapper: T::Hash[Symbol, T.any(Integer, T::Hash[Symbol, Integer])]
    ).returns(T::Array[T::Hash[Symbol, T.untyped]])
  end
  def build_custom_fields(data, mapper)
    mapper.each_with_object([]) do |(key, config), payload|
      value = data[key]
      next if value.blank?

      field_id, enum_id = config.is_a?(Hash) ? [ config[:field_id], config[:enum_id] ] : [ config, nil ]

      value_item = { value: value }
      value_item[:enum_id] = enum_id if enum_id

      payload << { field_id: field_id, values: [ value_item ] }
    end
  end

  sig { returns(String) }
  def note_text
    parts = []

    if @lead.survey_answers_data&.any?
      qa_text = @lead.survey_answers_data
        .map { "Вопрос: #{read_field(it, :question)}\nОтвет: #{read_field(it, :answer)}" }
        .join("\n\n")
      parts << qa_text
    end

    if @lead.courses_data&.any?
      courses_text = @lead.courses_data
        .map { "Курс: #{read_field(it, :slug)}\nПройдено уроков: #{read_field(it, :lessons_finished_count)}" }
        .join("\n\n")
      parts << "Информация по курсам:\n#{courses_text}"
    end

    parts.join("\n\n---\n\n")
  end

  sig { params(row: T.untyped, key: Symbol).returns(T.untyped) }
  def read_field(row, key)
    row.fetch(key)
  end
end
