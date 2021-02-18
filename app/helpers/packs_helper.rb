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
    javascript_pack_tag(*@javascript_packs, crossorigin: 'anonymous')
  end
end
# rubocop:enable Rails/HelperInstanceVariable
