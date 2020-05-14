# frozen_string_literal: true

class Web::LanguagesController < ApplicationController
  def show
    @language = Language.find(params[:id])
    @modules = @language.modules
  end
end
