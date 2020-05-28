class AddVersionToLesson < ActiveRecord::Migration[6.0]
  def change
    add_reference :language_module_lessons, :current_version, index: true
  end
end
