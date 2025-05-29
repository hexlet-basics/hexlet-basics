# typed: strict

class WorkflowLeadSerializer
  extend T::Sig

  LEAD_FIELD_MAPPER = {
    ym_client_id: 316941, # _ym_uid,
    utm_source: 316919,
    utm_medium: 316915,
    utm_campaign: 316917,
    utm_content: 316913,
    utm_term: 316921
  }.freeze

  # Соответствие ключей входных данных и ID полей AmoCRM для контакта
  CONTACT_FIELD_MAPPER = {
    user_id:   1_030_881,
    email:     { field_id: 316_907, enum_id: 158_463 }, # Личная почта
    phone:     { field_id: 316_905, enum_id: 158_453 }, # Мобильный телефон
    telegram:  925_081,
    whatsapp:  1_012_781
  }.freeze

  sig { params(event: LeadCreatedEvent).void }
  def initialize(event)
    @data = T.let(event.data, LeadCreatedEvent::DataShape)
  end

  sig { returns(T::Hash[Symbol, T.untyped]) }
  def to_h
    lead = base_payload

    lead_custom_fields = build_custom_fields(LEAD_FIELD_MAPPER)
    lead[:custom_fields_values] = lead_custom_fields unless lead_custom_fields.empty?

    lead[:_embedded] = { contacts: [ contact_payload ] }

    {
      lead: [ lead ],
      note: notes_payload
    }
  end

  private

  sig { returns(T::Hash[Symbol, T.untyped]) }
  def base_payload
    {
      name:                @data[:email],
      pipeline_id:         9_614_774,
      status_id:           76_748_378,
      responsible_user_id: 7_877_026
    }
  end

  sig { returns(T::Hash[Symbol, String]) }
  def notes_payload
    parts = []

    if @data[:survey_answers_data]&.any?
      qa_text = @data[:survey_answers_data]
                 .map { |qa| "Вопрос: #{qa[:question]}\nОтвет: #{qa[:answer]}" }
                 .join("\n\n")
      parts << qa_text
    end

    if @data[:courses_data]&.any?
      courses_text = @data[:courses_data]
                       .map { |cd| "Курс: #{cd[:slug]}\nПройдено уроков: #{cd[:lessons_finished_count]}" }
                       .join("\n\n")
      parts << "Информация по курсам:\n#{courses_text}"
    end

    { text: parts.join("\n\n---\n\n") }
  end

  sig { returns(T::Hash[Symbol, T.untyped]) }
  def contact_payload
    {
      name:                 @data[:user_name],
      first_name:           @data[:first_name],
      last_name:            @data[:last_name].to_s,
      custom_fields_values: build_custom_fields(CONTACT_FIELD_MAPPER)
    }
  end

  sig do
    params(
      mapper: T::Hash[Symbol, T.any(Integer, T::Hash[Symbol, Integer])]
    ).returns(T::Array[T::Hash[Symbol, T.untyped]])
  end
  def build_custom_fields(mapper)
    mapper.each_with_object([]) do |(key, config), arr|
      next unless @data[key]&.present?

      field_id, enum_id = if config.is_a?(Hash)
                            [ config[:field_id], config[:enum_id] ]
      else
                            [ config, nil ]
      end

      value_hash = { value: @data[key] }
      value_hash[:enum_id] = enum_id if enum_id

      arr << { field_id: field_id, values: [ value_hash ] }
    end
  end
end
