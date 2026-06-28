# typed: strict
# frozen_string_literal: true

class Web::Admin::Languages::ApplicationController < Web::Admin::ApplicationController
  helper_method :resource_language

  sig { returns(T.untyped) }
  def resource_language
    @resource_language ||= T.let(Language.find(params[:language_id]), T.untyped)
  end
end
