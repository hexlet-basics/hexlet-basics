# frozen_string_literal: true

class Web::Admin::Management::UsersController < Web::Admin::Management::ApplicationController
  include ActionController::Live

  def index
    query = params.fetch(:q, {}).with_defaults('s' => 'created_at desc')

    @search = User.includes(language_members: [ language: :current_version ]).ransack(query)

    users = @search.result

    respond_to do |format|
      format.html do
        @users = users.page(params[:page])
      end

      format.csv do
        response.headers['Last-Modified'] = Time.now.httpdate.to_s

        send_stream(filename: "users-#{params[:q][:language_members_created_at_gteq]}.csv") do |stream|
          stream.write "id, email, stack, finished_lessons\n"

          users.find_each do |user|
            language_members = user.language_members

            if language_members.any?
              language_members.each do |language_member|
                stream.write "#{user.id}, #{user.email}, #{language_member.language.name}, #{language_member.finished_lessons_count}\n"
              end
            else
              stream.write "#{user.id}, #{user.email}\n"
            end
          end
        end
      end
    end
  end

  def edit
    @user = User.find params[:id]
  end

  def update
    @user = User.find params[:id]

    if @user.update(user_params)
      redirect_to admin_management_users_path
    else
      render :edit
    end
  end

  private

  def user_params
    params.require(:user).permit(:admin)
  end
end
