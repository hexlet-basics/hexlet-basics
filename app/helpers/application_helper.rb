# frozen_string_literal: true

module ApplicationHelper
  include AuthConcern

  def let(value)
    yield value
  end

  def nav_menu_item(name, path = '#', options = {}, &block)
    assembled_options = options.merge(class: "nav-link text-dark #{options[:class]} #{active?(path)}".chomp)
    block_content = block ? capture(&block) : ''
    link = link_to name, path, assembled_options
    tag.li class: 'nav-item' do
      "#{link}#{block_content}".html_safe
    end
  end

  def active?(path, options = {})
    # raise options.inspect
    if options.key? :active_if
      'active' if options[:active_if]
    elsif current_page?(path)
      'active text-white'
    end
  end

  def structured_data_tag(path, args = {})
    content = render partial: "schemas/#{path}", formats: [:json], locals: args
    tag.script(content.html_safe, type: 'application/ld+json')
  end
end
