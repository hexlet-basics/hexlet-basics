# frozen_string_literal: true

module ApplicationHelper
  include AuthConcern

  def let(value)
    yield value
  end

  def nav_menu_item(name, path = '#', options = {})
    assembled_options = options.merge(class: "nav-link #{active?(path)}")
    tag.li class: 'nav-item' do
      link_to path, assembled_options do
        name
      end
    end
  end

  def active?(path, options = {})
    # raise options.inspect
    if options.key? :active_if
      'active' if options[:active_if]
    elsif current_page?(path)
      'active'
    end
  end

  def structured_data_tag(path, args = {})
    content = render partial: "schemas/#{path}", formats: [:json], locals: args
    tag.script(content.html_safe)
  end
end
