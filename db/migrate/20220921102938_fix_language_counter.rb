# frozen_string_literal: true

class FixLanguageCounter < ActiveRecord::Migration[7.0]
  def change
    Language::Member.counter_culture_fix_counts
    Language::Lesson.counter_culture_fix_counts
  end
end
