# typed: true

class Admin::LanguageCategoryForm < Language::Category
  include ActiveFormModel

  attr_writer :language_landing_page_ids

  permit :name, :slug, :header, :description,
    language_landing_page_ids: []

  def language_landing_page_ids
    @language_landing_page_ids || items.pluck(:language_landing_page_id)
  end

  def save(...)
    persisted = T.let(nil, T.untyped)

    transaction do
      persisted = super
      raise ActiveRecord::Rollback unless persisted

      sync_language_landing_page_ids!
    end

    persisted
  end

  def update(attributes)
    attrs = attributes.to_unsafe_h
    self.language_landing_page_ids = attrs["language_landing_page_ids"] || attrs[:language_landing_page_ids]

    updated = T.let(nil, T.untyped)

    transaction do
      updated = super(attrs.except("language_landing_page_ids", :language_landing_page_ids))
      raise ActiveRecord::Rollback unless updated

      sync_language_landing_page_ids!
    end

    updated
  end

  private

  def sync_language_landing_page_ids!
    normalized_ids = language_landing_page_ids.filter_map do |id|
      next if id.blank?

      id.to_i
    end.uniq

    items.where.not(language_landing_page_id: normalized_ids).delete_all

    existing_ids = items.where(language_landing_page_id: normalized_ids).pluck(:language_landing_page_id)
    missing_ids = normalized_ids - existing_ids

    missing_ids.each do |landing_page_id|
      items.create!(language_landing_page_id: landing_page_id)
    end
  end
end
