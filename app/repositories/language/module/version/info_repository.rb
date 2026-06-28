# typed: strict
# frozen_string_literal: true

module Language::Module::Version::InfoRepository
  extend ActiveSupport::Concern
  extend T::Helpers
  requires_ancestor { ActiveRecord::Base }
  include LocaleRepository

  included do
    T.bind(self, T.class_of(ActiveRecord::Base))
  end
end
