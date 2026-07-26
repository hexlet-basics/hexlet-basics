# typed: strict
# frozen_string_literal: true

class Ahoy::Store < Ahoy::DatabaseStore
  extend T::Sig

  sig { params(data: T.untyped).returns(T.untyped) }
  def track_visit(data)
    ym_client_id = extract_ym_client_id(data)

    payload = data.dup
    payload[:ym_client_id] ||= ym_client_id if ym_client_id

    super(payload)
  end

  sig { params(data: T.untyped).returns(T.untyped) }
  def track_event(data)
    super(data)

    ym_client_id = extract_ym_client_id(data)
    return unless ym_client_id

    current_visit = visit
    return if current_visit.ym_client_id?

    current_visit.update!(ym_client_id:)
  end

  sig { params(data: T.untyped).returns(T.untyped) }
  def geocode(data)
    super
  end

  sig { params(data: T.untyped).returns(T.untyped) }
  def authenticate(data)
    super
  end

  private

  sig { params(data: T.untyped).returns(T.untyped) }
  def extract_ym_client_id(data)
    value = data.dig(:properties, :ym_client_id) ||
      data.dig(:properties, "ym_client_id") ||
      request&.params&.dig("ym_client_id")

    value = value.to_s.strip
    value.presence
  end
end
