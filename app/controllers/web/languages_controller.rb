# frozen_string_literal: true

class Web::LanguagesController < Web::ApplicationController
  def show
    @language = Language.find(params[:id])
    @modules = @language.modules
  end
end
