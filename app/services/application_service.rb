# typed: strict
# frozen_string_literal: true

class ApplicationService
  extend T::Sig

  private

  sig { returns(Typed::Success[NilClass]) }
  def self.success
    Typed::Success.blank
  end

  sig do
    type_parameters(:Payload)
      .params(payload: T.type_parameter(:Payload))
      .returns(Typed::Success[T.type_parameter(:Payload)])
  end
  def self.success_with(payload)
    Typed::Success.new(payload)
  end

  sig { returns(Typed::Failure[NilClass]) }
  def self.fail
    Typed::Failure.blank
  end

  sig do
    type_parameters(:Error)
      .params(error: T.type_parameter(:Error))
      .returns(Typed::Failure[T.type_parameter(:Error)])
  end
  def self.fail_with(error)
    Typed::Failure.new(error)
  end
end
