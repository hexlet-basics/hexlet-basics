class AddCoursesDataToLeads < ActiveRecord::Migration[8.0]
  def change
    add_column :leads, :courses_data, :text
  end
end
