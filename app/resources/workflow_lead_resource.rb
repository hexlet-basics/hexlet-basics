# typed: strict

class WorkflowLeadResource < ApplicationResource
  extend T::Sig

  typelize_from Lead

  attributes :lead, :note

  FormData = T.type_alias { T::Hash[Symbol, T.untyped] }
  NotePayload = T.type_alias { { text: String } }

  FieldMeta = T.type_alias { T::Hash[Symbol, T.untyped] }

  LEAD_FIELD_MAPPER = T.let({
    utm_content: { id: 316_913, code: "UTM_CONTENT" },
    utm_medium: { id: 316_915, code: "UTM_MEDIUM" },
    utm_campaign: { id: 316_917, code: "UTM_CAMPAIGN" },
    utm_source: { id: 316_919, code: "UTM_SOURCE" },
    utm_term: { id: 316_921, code: "UTM_TERM" },
    utm_referrer: { id: 316_923, code: "UTM_REFERRER" },
    roistat: { id: 316_925, code: "ROISTAT" },
    referrer: { id: 316_927, code: "REFERRER" },
    openstat_service: { id: 316_929, code: "OPENSTAT_SERVICE" },
    openstat_campaign: { id: 316_931, code: "OPENSTAT_CAMPAIGN" },
    openstat_ad: { id: 316_933, code: "OPENSTAT_AD" },
    openstat_source: { id: 316_935, code: "OPENSTAT_SOURCE" },
    from: { id: 316_937, code: "FROM" },
    gclientid: { id: 316_939, code: "GCLIENTID" },
    ym_client_id: { id: 316_941, code: "_YM_UID" },
    ym_counter: { id: 316_943, code: "_YM_COUNTER" }, # from ENV["YM_COUNTER"]
    gclid: { id: 316_945, code: "GCLID" },
    yclid: { id: 316_947, code: "YCLID" },
    fbclid: { id: 316_949, code: "FBCLID" },
    ga_utm: { id: 957_711, code: "GA_UTM" },
    telegram: 925_081,
    whatsapp: 1_012_781,
    source_form: { id: 936_587, code: nil }
  }.freeze, T::Hash[Symbol, FieldMeta])

  sig { params(_lead: Lead).returns(T::Array[FormData]) }
  def lead(_lead)
    source_form = "lead_form" # object.source_form.presence || "lead_form"
    lead_entity = lead_payload(source_form)
    lead_custom_fields = build_lead_custom_fields(source_form)
    lead_entity[:custom_fields_values] = lead_custom_fields unless lead_custom_fields.empty?

    [
      {
        source_uid: source_form,
        source_name: source_form,
        metadata: metadata_payload(source_form),
        _embedded: {
          contacts: [ contact_payload ],
          leads: [ lead_entity ]
        }
      }
    ]
  end

  sig { params(_lead: Lead).returns(NotePayload) }
  def note(_lead)
    notes_payload
  end

  private

  sig { params(source_form: String).returns(T::Hash[Symbol, T.untyped]) }
  def lead_payload(source_form)
    {
      name: object.email || "Lead from #{source_form}",
      pipeline_id: 9_614_774,
      # status_id: 76_748_378,
      responsible_user_id: 7_877_026
    }
  end

  sig { params(source_form: String).returns(T::Hash[Symbol, T.untyped]) }
  def metadata_payload(source_form)
    # AmoCRM does not accept IPv6.
    visit = object.ahoy_visit
    is_ipv4 = IPAddr.new(visit&.ip.to_s).ipv4? rescue false

    {
      form_id: source_form,
      form_name: source_form,
      form_page: visit&.landing_page,
      form_sent_at: object.created_at.to_i,
      ip: is_ipv4 ? visit&.ip : nil,
      referer: visit&.referrer,
      visitor_uid: visit&.visitor_token
    }.compact
  end

  sig { returns(NotePayload) }
  def notes_payload
    parts = []

    if object.survey_answers_data&.any?
      qa_text = object.survey_answers_data
        .map { |qa| "Вопрос: #{qa[:question]}\nОтвет: #{qa[:answer]}" }
        .join("\n\n")
      parts << qa_text
    end

    if object.courses_data&.any?
      courses_text = object.courses_data
        .map { |cd| "Курс: #{cd[:slug]}\nПройдено уроков: #{cd[:lessons_finished_count]}" }
        .join("\n\n")
      parts << "Информация по курсам:\n#{courses_text}"
    end

    # if object.additional_info.present?
    #   parts << "Дополнительная информация:\n#{object.additional_info}"
    # end

    { text: parts.join("\n\n---\n\n") }
  end

  sig { returns(T::Hash[Symbol, T.untyped]) }
  def contact_payload
    user = object.user
    email = object.email || user&.email

    {
      name: contact_name(user, email),
      first_name: user&.first_name.to_s,
      last_name: user&.last_name.to_s,
      custom_fields_values: build_contact_fields(email, object.phone_number)
    }
  end

  sig { params(user: T.nilable(User), email: T.nilable(String)).returns(String) }
  def contact_name(user, email)
    user&.first_name.presence ||
    email.presence ||
    object.phone.presence ||
    object.telegram.presence ||
    object.whatsapp.presence ||
    "Unknown"
  end

  sig do
    params(email: T.nilable(String), phone: T.nilable(String)).returns(T::Array[T::Hash[Symbol, T.untyped]])
  end
  def build_contact_fields(email, phone)
    fields = []
    if email.present?
      fields << {
        field_code: "EMAIL",
        values: [ { value: email, enum_code: "WORK" } ]
      }
    end

    if phone.present?
      fields << {
        field_code: "PHONE",
        values: [ { value: phone, enum_code: "WORK" } ]
      }
    end

    fields
  end

  sig { params(source_form: String).returns(T::Array[T::Hash[Symbol, T.untyped]]) }
  def build_lead_custom_fields(source_form)
    visit = object.ahoy_visit
    tracking_values = extract_tracking_values

    data = {
      utm_source: visit&.utm_source,
      utm_medium: visit&.utm_medium,
      utm_campaign: visit&.utm_campaign,
      utm_content: visit&.utm_content,
      utm_term: visit&.utm_term,
      utm_referrer: visit&.referrer,
      # roistat: tracking_values[:roistat],
      referrer: visit&.referrer,
      # openstat_service: tracking_values[:openstat_service],
      # openstat_campaign: tracking_values[:openstat_campaign],
      # openstat_ad: tracking_values[:openstat_ad],
      # openstat_source: tracking_values[:openstat_source],
      from: tracking_values[:from],
      # gclientid: object.ga_client_id || tracking_values[:gclientid],
      ym_client_id: visit&.ym_client_id,
      ym_counter: ENV.fetch("YM_COUNTER"),
      gclid: tracking_values[:gclid],
      yclid: object.ym_client_id || tracking_values[:yclid],
      fbclid: tracking_values[:fbclid],
      ga_utm: tracking_values[:ga_utm],
      source_form: source_form
    }

    data.each_with_object([]) do |(key, value), result|
      next unless value.present?

      field_meta = T.let(LEAD_FIELD_MAPPER.fetch(key), FieldMeta)
      payload = {
        field_id: field_meta.fetch(:id),
        values: [ { value: value } ]
      }
      field_code = field_meta[:code]
      payload[:field_code] = field_code if field_code.present?
      result << payload
    end
  end

  sig { returns(T::Hash[Symbol, T.nilable(String)]) }
  def extract_tracking_values
    urls = [ object.ahoy_visit&.landing_page ].compact
    merged = urls.each_with_object({}) do |url, acc|
      acc.merge!(extract_query_params(url))
    end

    {
      # roistat: merged["roistat"],
      # openstat_service: merged["openstat_service"],
      # openstat_campaign: merged["openstat_campaign"],
      # openstat_ad: merged["openstat_ad"],
      # openstat_source: merged["openstat_source"],
      from: merged["from"],
      gclientid: merged["gclientid"],
      gclid: merged["gclid"],
      yclid: merged["yclid"],
      fbclid: merged["fbclid"],
      ga_utm: merged["ga_utm"]
    }
  end

  sig { params(url: T.nilable(String)).returns(T::Hash[String, String]) }
  def extract_query_params(url)
    return {} if url.blank?

    uri = URI.parse(url)
    return {} if uri.query.blank?

    Rack::Utils.parse_nested_query(uri.query).transform_values(&:to_s)
  rescue URI::InvalidURIError
    {}
  end
end
