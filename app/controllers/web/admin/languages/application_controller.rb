# typed: strict
# frozen_string_literal: true

class Web::Admin::Languages::ApplicationController < Web::Admin::ApplicationController
  helper_method :resource_language

  sig { returns(Language) }
  def resource_language
    @resource_language ||= T.let(Language.find(params[:language_id]), T.nilable(Language))
  end
end
