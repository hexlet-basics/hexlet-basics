# frozen_string_literal: true

class Web::PagesController < Web::ApplicationController
  def show
    title params[:id].to_sym
    render params[:id]
  end
end
