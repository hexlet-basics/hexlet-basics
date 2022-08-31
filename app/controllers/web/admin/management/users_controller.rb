# frozen_string_literal: true

class Web::Admin::Management::UsersController < Web::Admin::Management::ApplicationController
  def index
    query = params.fetch(:q, {}).with_defaults('s' => 'created_at desc')

    @search = User.includes(language_members: :language).ransack(query)

    users = @search.result

    respond_to do |format|
      format.html do
        @users = users.page(params[:page])
      end

      format.csv do
        fields = %w[id email stack finished_lessons]

        csv_string = CSV.generate do |csv|
          csv << fields

          users.find_each do |user|
            language_members = user.language_members

            if language_members.any?
              language_members.each do |language_member|
                csv << [
                  user.id,
                  user.email,
                  language_member.language.name,
                  language_member.finished_lessons_count
                ]
              end
            else
              csv << [
                user.id,
                user.email
              ]
            end
          end
        end

        send_data csv_string, filename: "users-#{Time.zone.today}.csv"
      end
    end
  end
end
