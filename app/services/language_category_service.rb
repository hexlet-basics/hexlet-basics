# typed: strict
# frozen_string_literal: true

class LanguageCategoryService < ApplicationService
  Category = Language::Category

  class << self
    extend T::Sig

    sig { params(struct: LanguageCategoryStruct, locale: String).returns(Typed::Result[Category, Category]) }
    def create(struct, locale:)
      category = Category.new(category_attributes(struct).merge(locale:))
      persist(category, struct)
    end

    sig { params(id: String, struct: LanguageCategoryStruct, locale: String).returns(Typed::Result[Category, Category]) }
    def update(id, struct, locale:)
      category = Category.find(id)
      category.assign_attributes(category_attributes(struct).merge(locale:))
      persist(category, struct)
    end

    private

    sig { params(category: Category, struct: LanguageCategoryStruct).returns(Typed::Result[Category, Category]) }
    def persist(category, struct)
      saved = T.let(false, T::Boolean)
      ActiveRecord::Base.transaction do
        saved = category.save
        raise ActiveRecord::Rollback unless saved

        sync_items!(category, struct.language_landing_page_ids || [])
      end

      saved ? success_with(category) : fail_with(category)
    end

    # language_landing_page_ids isn't a column — it maps to join records
    sig { params(struct: LanguageCategoryStruct).returns(T::Hash[Symbol, T.untyped]) }
    def category_attributes(struct)
      struct.attributes.except(:language_landing_page_ids)
    end

    sig { params(category: Category, ids: T::Array[Integer]).void }
    def sync_items!(category, ids)
      normalized = ids.uniq
      category.items.where.not(language_landing_page_id: normalized).delete_all
      existing = category.items.where(language_landing_page_id: normalized).pluck(:language_landing_page_id)
      (normalized - existing).each do |landing_page_id|
        category.items.create!(language_landing_page_id: landing_page_id)
      end
    end
  end
end
