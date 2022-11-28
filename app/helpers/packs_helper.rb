# frozen_string_literal: true

# rubocop:disable Rails/HelperInstanceVariable
module PacksHelper
  def append_javascript_packs(*packs)
    @javascript_packs ||= []
    @javascript_packs.push(*packs)
  end

  def prepend_javascript_packs(*packs)
    @javascript_packs ||= []
    @javascript_packs.unshift(*packs)
  end

  def render_javascript_stored_packs
    javascript_include_tag(*@javascript_packs&.uniq, crossorigin: 'anonymous', defer: true, nonce: true)
  end

  def append_stylesheet_packs(*packs)
    @stylesheet_packs ||= []
    @stylesheet_packs.push(*packs)
  end

  def render_stylesheet_stored_packs
    stylesheet_link_tag(*@stylesheet_packs&.uniq, media: 'all')
  end
end
# rubocop:enable Rails/HelperInstanceVariable
