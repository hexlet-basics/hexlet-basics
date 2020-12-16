# frozen_string_literal: true

module TitleConcern
  extend ActiveSupport::Concern

  included do
    helper_method :title
  end

  def base_title
    title t('title.base')
  end

  def title(part = nil, options = {})
    @parts ||= []

    unless part
      return nil if @parts.blank?

      return @parts.reject(&:blank?).reverse.join(' - ').truncate(65, separator: ' ')
    end

    @parts << if part.instance_of? Symbol
                translate_title(part, options)
              else
                part
              end
  end

  def translate_title(key, options = {})
    scope = [:titles]
    if options[:scope]
      scope += options[:scope].split('.') if options[:scope].is_a? String
      scope += options[:scope] if options[:scope].is_a? Array
    else
      scope << params[:controller].split('/')
    end
    # scope << params[:action]

    # values = options[:values] || {}

    # msg = I18n.t(key, scope: scope, **values)
    # #NOTE color logging
    # Rails.logger.debug(Term::ANSIColor.green("title: #{msg}"))

    I18n.t(key, scope: scope)

    # msg
  end
end
