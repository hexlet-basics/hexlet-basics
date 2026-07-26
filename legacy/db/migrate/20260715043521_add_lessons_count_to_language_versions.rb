class AddLessonsCountToLanguageVersions < ActiveRecord::Migration[8.1]
  def change
    # counter_culture column for Language::Lesson::Version.counter_culture :language_version.
    # Column only — existing rows are backfilled out-of-band via
    # Language::Lesson::Version.counter_culture_fix_counts (not in the migration).
    add_column :language_versions, :lessons_count, :integer, default: 0, null: false
  end
end
