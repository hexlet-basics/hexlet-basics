# frozen_string_literal: true

class Web::PagesController < Web::ApplicationController
  def show
    set_meta_tags title: t(params[:id], scope: 'web.pages')

    render params[:id]
  end
end
