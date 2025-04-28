namespace :analytics do
  desc "Refresh user survey pivot table"
  task refresh_user_survey_pivot: :environment do
    sql_file_path = Rails.root.join("db/sqls/user_survey_pivot.sql")

    unless File.exist?(sql_file_path)
      raise "SQL file not found: #{sql_file_path}"
    end

    sql = File.read(sql_file_path)

    ActiveRecord::Base.transaction do
      puts "[UserSurveyPivot] Running SQL from #{sql_file_path}..."
      ActiveRecord::Base.connection.execute(sql)
      puts "[UserSurveyPivot] Refresh complete!"
    end
  end
end
