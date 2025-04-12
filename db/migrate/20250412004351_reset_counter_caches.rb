class ResetCounterCaches < ActiveRecord::Migration[8.0]
  def change
    Language::Lesson::Member::Message.counter_culture_fix_counts
  end
end
